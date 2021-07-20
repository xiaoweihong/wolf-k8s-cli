package cmd

import (
	"github.com/spf13/cobra"
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
	"log"
	config "wolf-k8s-cli/configs"
	"wolf-k8s-cli/k8s"

	"context"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var monadoCmd = &cobra.Command{
	Use:   "monado",
	Short: "修改monado存图时间",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := k8s.NewClient(*config.Kubeconfig, nil)
		if err != nil {
			log.Println(err)
			return
		}
		cms, err := clientset.CoreV1().ConfigMaps("wolf").List(context.Background(), metav1.ListOptions{
			FieldSelector: "metadata.name=monado-capture-config"})
		if err != nil {
			log.Panicln(err)
		}

		monadoCapturedCmConfigJson := cms.Items[0].Data["config.json"]
		monadoCapturedCmEntrypointSh := cms.Items[0].Data["entrypoint.sh"]
		log.Println("原图片存储时间:", gjson.Get(monadoCapturedCmConfigJson, "StorageConfigs").Get("node3").Get("PersistenceConfig").Get("ConfigInfo").Get("TTL"))
		value, err := sjson.Set(monadoCapturedCmConfigJson, "StorageConfigs.node3.PersistenceConfig.ConfigInfo.TTL", config.TTL)
		if err != nil {
			log.Panicln(err)
		}
		v1.ConfigMap("monado-capture-config", "wolf")
		res, err := clientset.CoreV1().ConfigMaps("wolf").Update(context.Background(), &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: "monado-capture-config",
			},
			Data: map[string]string{
				"config.json":   value,
				"entrypoint.sh": monadoCapturedCmEntrypointSh,
			},
		}, metav1.UpdateOptions{})

		if err != nil {
			log.Panicln(err)
		}
		log.Println("修改成功，新图片存储时间:", gjson.Get(res.Data["config.json"], "StorageConfigs").Get("node3").Get("PersistenceConfig").Get("ConfigInfo").Get("TTL"))
		log.Println("开始重启pod，使配置生效...")
		//list, err := clientset.AppsV1().Deployments("wolf").Patch(context.Background(),"monado","", nil,metav1.PatchOptions{})
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(list)
		//fmt.Println(list.Items[0].Spec)
		//clientset.CoreV1().Pods(*namespace).Delete(context.Background(),"",metav1.DeleteOptions{})

		list, err := clientset.CoreV1().Pods("wolf").List(context.Background(), metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=monado",
		})
		if err != nil {
			log.Panicln(err)
		}
		log.Println("原pod名:", list.Items[0].Name)

		err = clientset.CoreV1().Pods("wolf").Delete(context.Background(), list.Items[0].Name, metav1.DeleteOptions{})
		if err != nil {
			panic(err)
		}
		log.Println("pod重启成功...")
		log.Println(`修改完毕...等待pod重启完毕(kubectl get pods -n wolf |grep monado)`)
	},
}

func init() {
	monadoCmd.Flags().StringVar(&config.TTL, "ttl", "6M", "图片存储时间")
}
