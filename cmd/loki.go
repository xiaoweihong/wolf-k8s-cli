package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	config "wolf-k8s-cli/configs"
	"wolf-k8s-cli/k8s"
)

var lokiCmd = &cobra.Command{
	Use:   "loki",
	Short: "修改loki的一些常用配置",
}

var deepcloudSwitchCmd = &cobra.Command{
	Use:   "deepcloud",
	Short: "调整深瞳云同步开关",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := k8s.NewClient(*config.Kubeconfig, nil)
		if err != nil {
			log.Println(err)
			return
		}
		configMapClient := clientset.CoreV1().ConfigMaps("wolf")
		lokiConfigs, err := configMapClient.Get(context.Background(), "loki-config", metav1.GetOptions{})
		if err != nil {
			log.Panicln(err)
		}
		defaultJson := lokiConfigs.Data["default.json"]
		productionJson := lokiConfigs.Data["production.json"]
		switchBool := gjson.Get(defaultJson, "loki").Get("base").Get("deepcloud").Get("switch").Bool()
		log.Println("修改前深瞳云开关为:", switchBool)
		switchChange, err := sjson.Set(defaultJson, "loki.base.deepcloud.switch", !switchBool)
		if err != nil {
			log.Panicln(err)
		}
		var str bytes.Buffer
		json.Indent(&str, []byte(switchChange), "", " ")
		res, err := configMapClient.Update(context.Background(), &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: "loki-config",
			},
			Data: map[string]string{
				"default.json":    str.String(),
				"production.json": productionJson,
			},
		}, metav1.UpdateOptions{})
		if err != nil {
			log.Panicln(err)
		}
		switchBool = gjson.Get(res.Data["default.json"], "loki").Get("base").Get("deepcloud").Get("switch").Bool()
		log.Println("修改后的深瞳云参数开关为:", switchBool)
		if !switchBool {
			log.Println("loki可以删除前端无法删除的设备")
			log.Println("删除后记得再次执行改命令改回,否则会导致后面同步或者新建的设备无数据")
		} else {
			log.Println("配置正常")
		}
		list, err := clientset.CoreV1().Pods("wolf").List(context.Background(), metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=loki",
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
		log.Println(`修改完毕...等待pod重启完毕(kubectl get pods -n wolf |grep loki)`)
	},
}

func init() {
	rootCmd.AddCommand(lokiCmd)

	lokiCmd.AddCommand(deepcloudSwitchCmd)
}
