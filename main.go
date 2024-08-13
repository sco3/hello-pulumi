package main

import (
	"fmt"
	_ "fmt"
	"os"

	appv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	const softExcludesTxt = "soft-excludes.txt"
	const appName = "cortex"
	const configMapName = pulumi.String("my-config-map")
	const configVolume = "config-volume"

	fmt.Printf(".start.\n")

	pulumi.Run(func(ctx *pulumi.Context) error {
		appLabels := pulumi.StringMap{
			"app": pulumi.String(appName),
		}
		fmt.Printf("labels: %v\n", appLabels)
		softExcludes, err := os.ReadFile(softExcludesTxt)
		if err != nil {
			return err
		}
		fmt.Printf("excludes:\n%v\n", string(softExcludes))
		cfgMap, err := corev1.NewConfigMap(ctx, appName, &corev1.ConfigMapArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: appLabels,
				Name:   configMapName,
			},
			Data: pulumi.StringMap{
				softExcludesTxt: pulumi.String(softExcludes),
			},
		})
		if err != nil {
			return err
		}
		fmt.Printf("config map:\n%v\n", cfgMap)

		const statefulSetName = "my-statefulset"
		statefulSet, err := appv1.NewStatefulSet(ctx, statefulSetName, &appv1.StatefulSetArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String(statefulSetName),
			},
			Spec: &appv1.StatefulSetSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app": pulumi.String(appName),
					},
				},
				ServiceName: pulumi.String("my-service"),
				Replicas:    pulumi.Int(1),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String(appName),
						},
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String(appName),
								Image: pulumi.String("localhost/cortex:TTDP-4489"),
								VolumeMounts: corev1.VolumeMountArray{
									&corev1.VolumeMountArgs{
										Name:      pulumi.String(configVolume),
										MountPath: pulumi.String("/home/skyhook-user/"),
										ReadOnly:  pulumi.Bool(false),
									},
								},
							},
						},
						Volumes: corev1.VolumeArray{
							&corev1.VolumeArgs{
								Name: pulumi.String(configVolume),
								ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
									Name: configMapName,
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err

		}
		fmt.Errorf("StatefulSet: %v", statefulSet)
		return nil
	})

	fmt.Printf(".finish.\n")
}
