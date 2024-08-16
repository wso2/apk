/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package proxy

import (
	"fmt"
	"strings"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/bootstrap"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/envoy"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/infrastructure/kubernetes/resource"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/ir"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/version"

	egv1a1 "github.com/wso2/apk/adapter/internal/operator/gateway-api/v1alpha1"
	utils "github.com/wso2/apk/adapter/pkg/utils/misc"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
)

const (
	SdsCAFilename   = "xds-trusted-ca.json"
	SdsCertFilename = "xds-certificate.json"
	// XdsTLSCertFilename is the fully qualified path of the file containing Envoy's
	// xDS server TLS certificate.
	XdsTLSCertFilename = "/home/wso2/security/keystore/router.crt"
	// XdsTLSCertFilename = "/home/wso2/security/keystore/router.crt"
	// XdsTLSKeyFilename is the fully qualified path of the file containing Envoy's
	// xDS server TLS key.
	XdsTLSKeyFilename = "/home/wso2/security/keystore/router.key"
	// XdsTLSKeyFilename = "/home/wso2/security/keystore/router.key"
	// XdsTLSCaFilename is the fully qualified path of the file containing Envoy's
	// trusted CA certificate.
	// XdsTLSCaFilename = "/certs/ca.crt"
	XdsTLSCaFilename = "/home/wso2/security/truststore/adapter.crt"
	// envoyContainerName is the name of the Envoy container.
	envoyContainerName    = "envoy"
	enforcerContainerName = "enforcer"
	// envoyNsEnvVar is the name of the APK Gateway namespace environment variable.
	envoyNsEnvVar = "ENVOY_GATEWAY_NAMESPACE"
	// envoyPodEnvVar is the name of the Envoy pod name environment variable.
	envoyPodEnvVar = "ENVOY_POD_NAME"
)

var (
	// xDS certificate rotation is supported by using SDS path-based resource files.
	SdsCAConfigMapData = fmt.Sprintf(`{"resources":[{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.Secret",`+
		`"name":"xds_trusted_ca","validation_context":{"trusted_ca":{"filename":"%s"},`+
		`"match_typed_subject_alt_names":[{"san_type":"DNS","matcher":{"exact":"envoy-gateway"}}]}}]}`, XdsTLSCaFilename)
	SdsCertConfigMapData = fmt.Sprintf(`{"resources":[{"@type":"type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.Secret",`+
		`"name":"xds_certificate","tls_certificate":{"certificate_chain":{"filename":"%s"},`+
		`"private_key":{"filename":"%s"}}}]}`, XdsTLSCertFilename, XdsTLSKeyFilename)
)

// ExpectedResourceHashedName returns expected resource hashed name including up to the 48 characters of the original name.
func ExpectedResourceHashedName(name string) string {
	hashedName := utils.GetHashedName(name, 35)
	return hashedName
}

// EnvoyAppLabel returns the labels used for all Envoy resources.
func EnvoyAppLabel() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "envoy",
		"app.kubernetes.io/component":  "proxy",
		"app.kubernetes.io/managed-by": "apk-gateway",
	}
}

// EnvoyAppLabelSelector returns the labels used for all Envoy resources.
func EnvoyAppLabelSelector() []string {
	return []string{
		"app.kubernetes.io/name=envoy",
		"app.kubernetes.io/component=proxy",
		"app.kubernetes.io/managed-by=apk-gateway",
	}
}

// envoyLabels returns the labels, including extraLabels, used for Envoy resources.
func envoyLabels(extraLabels map[string]string) map[string]string {
	labels := EnvoyAppLabel()
	for k, v := range extraLabels {
		labels[k] = v
	}

	return labels
}

func enablePrometheus(infra *ir.ProxyInfra) bool {
	return false
}

// expectedProxyContainers returns expected proxy containers.
func expectedProxyContainers(infra *ir.ProxyInfra,
	deploymentConfig *egv1a1.KubernetesDeploymentSpec,
	shutdownConfig *egv1a1.ShutdownConfig) ([]corev1.Container, error) {
	// Define slice to hold container ports
	var ports []corev1.ContainerPort

	// Iterate over listeners and ports to get container ports
	for _, listener := range infra.Listeners {
		for _, p := range listener.Ports {
			var protocol corev1.Protocol
			switch p.Protocol {
			case ir.HTTPProtocolType, ir.HTTPSProtocolType, ir.TLSProtocolType, ir.TCPProtocolType:
				protocol = corev1.ProtocolTCP
			case ir.UDPProtocolType:
				protocol = corev1.ProtocolUDP
			default:
				return nil, fmt.Errorf("invalid protocol %q", p.Protocol)
			}
			port := corev1.ContainerPort{
				// hashed container port name including up to the 6 characters of the port name and the maximum of 15 characters.
				Name:          utils.GetHashedName(p.Name, 6),
				ContainerPort: p.ContainerPort,
				Protocol:      protocol,
			}
			ports = append(ports, port)
		}
	}
	port := corev1.ContainerPort{
		// hashed container port name including up to the 6 characters of the port name and the maximum of 15 characters.
		Name:          "admin",
		ContainerPort: 9000,
		Protocol:      "TCP",
	}
	ports = append(ports, port)
	if enablePrometheus(infra) {
		ports = append(ports, corev1.ContainerPort{
			Name:          "metrics",
			ContainerPort: bootstrap.EnvoyReadinessPort, // TODO: make this configurable
			Protocol:      corev1.ProtocolTCP,
		})
	}

	var bootstrapConfigurations string

	var proxyMetrics *egv1a1.ProxyMetrics

	// Get the default Bootstrap
	bootstrapConfigurations, err := bootstrap.GetRenderedBootstrapConfig(proxyMetrics)
	if err != nil {
		return nil, err
	}

	args := []string{
		fmt.Sprintf("--service-cluster %s", infra.Name),
		fmt.Sprintf("--service-node $(%s)", envoyPodEnvVar),
		fmt.Sprintf("--config-yaml %s", bootstrapConfigurations),
		fmt.Sprintf("--log-level %s", "warn"),
		"--cpuset-threads",
	}

	// if infra.Config != nil &&
	// 	infra.Config.Spec.Concurrency != nil {
	// 	args = append(args, fmt.Sprintf("--concurrency %d", *infra.Config.Spec.Concurrency))
	// }

	// if componentsLogLevel := logging.GetEnvoyProxyComponentLevel(); componentsLogLevel != "" {
	// 	args = append(args, fmt.Sprintf("--component-log-level %s", componentsLogLevel))
	// }

	if shutdownConfig != nil && shutdownConfig.DrainTimeout != nil {
		args = append(args, fmt.Sprintf("--drain-time-s %.0f", shutdownConfig.DrainTimeout.Seconds()))
	}

	// if infra.Config != nil {
	// 	args = append(args, infra.Config.Spec.ExtraArgs...)
	// }

	containers := []corev1.Container{
		{
			Name:                     envoyContainerName,
			Image:                    *deploymentConfig.EnvoyProxyContainer.Image,
			ImagePullPolicy:          corev1.PullIfNotPresent,
			Command:                  []string{"envoy"},
			Args:                     args,
			Env:                      expectedContainerEnv(deploymentConfig.EnvoyProxyContainer),
			Resources:                *deploymentConfig.EnvoyProxyContainer.Resources,
			SecurityContext:          deploymentConfig.EnvoyProxyContainer.SecurityContext,
			Ports:                    ports,
			VolumeMounts:             expectedContainerVolumeMounts(deploymentConfig.EnvoyProxyContainer),
			TerminationMessagePolicy: corev1.TerminationMessageReadFile,
			TerminationMessagePath:   "/dev/termination-log",
			ReadinessProbe: &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path:   bootstrap.EnvoyReadinessPath,
						Port:   intstr.IntOrString{Type: intstr.Int, IntVal: bootstrap.EnvoyReadinessPort},
						Scheme: corev1.URISchemeHTTP,
					},
				},
				TimeoutSeconds:   1,
				PeriodSeconds:    10,
				SuccessThreshold: 1,
				FailureThreshold: 3,
			},
			Lifecycle: &corev1.Lifecycle{
				PreStop: &corev1.LifecycleHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path:   envoy.ShutdownManagerReadyPath,
						Port:   intstr.FromInt32(envoy.ShutdownManagerPort),
						Scheme: corev1.URISchemeHTTP,
					},
				},
			},
		},
		{
			Name:                     "shutdown-manager",
			Image:                    expectedShutdownManagerImage(),
			ImagePullPolicy:          corev1.PullIfNotPresent,
			Command:                  []string{"envoy-gateway"},
			Args:                     expectedShutdownManagerArgs(shutdownConfig),
			Env:                      expectedContainerEnv(nil),
			Resources:                *egv1a1.DefaultShutdownManagerContainerResourceRequirements(),
			TerminationMessagePolicy: corev1.TerminationMessageReadFile,
			TerminationMessagePath:   "/dev/termination-log",
			ReadinessProbe: &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path:   envoy.ShutdownManagerHealthCheckPath,
						Port:   intstr.IntOrString{Type: intstr.Int, IntVal: envoy.ShutdownManagerPort},
						Scheme: corev1.URISchemeHTTP,
					},
				},
				TimeoutSeconds:   1,
				PeriodSeconds:    10,
				SuccessThreshold: 1,
				FailureThreshold: 3,
			},
			LivenessProbe: &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path:   envoy.ShutdownManagerHealthCheckPath,
						Port:   intstr.IntOrString{Type: intstr.Int, IntVal: envoy.ShutdownManagerPort},
						Scheme: corev1.URISchemeHTTP,
					},
				},
				TimeoutSeconds:   1,
				PeriodSeconds:    10,
				SuccessThreshold: 1,
				FailureThreshold: 3,
			},
			Lifecycle: &corev1.Lifecycle{
				PreStop: &corev1.LifecycleHandler{
					Exec: &corev1.ExecAction{
						Command: expectedShutdownPreStopCommand(shutdownConfig),
					},
				},
			},
		},
		{
			Name:            enforcerContainerName,
			Image:           *deploymentConfig.EnforcerContainer.Image,
			ImagePullPolicy: corev1.PullIfNotPresent,
			// Command:                  []string{"envoy"},
			// Args:                     args,
			Env:                      expectedEnforcerEnv(deploymentConfig.EnforcerContainer),
			Resources:                *deploymentConfig.EnforcerContainer.Resources,
			SecurityContext:          deploymentConfig.EnforcerContainer.SecurityContext,
			Ports:                    expectedEnforcerPorts(),
			VolumeMounts:             expectedEnforcerVolumeMounts(deploymentConfig.EnforcerContainer),
			TerminationMessagePolicy: corev1.TerminationMessageReadFile,
			TerminationMessagePath:   "/dev/termination-log",
			//todo(amali)
			// 	ReadinessProbe: &corev1.Probe{
			// 		ProbeHandler: corev1.ProbeHandler{
			// 			Exec: &corev1.ExecAction{
			// 				Command: []string{"sh", "check_health.sh"},
			// 			},
			// 		},
			// 		TimeoutSeconds:      1,
			// 		PeriodSeconds:       20,
			// 		SuccessThreshold:    1,
			// 		FailureThreshold:    5,
			// 		InitialDelaySeconds: 20,
			// 	},
			// 	LivenessProbe: &corev1.Probe{
			// 		ProbeHandler: corev1.ProbeHandler{
			// 			Exec: &corev1.ExecAction{
			// 				Command: []string{"sh", "check_health.sh"},
			// 			},
			// 		},
			// 		TimeoutSeconds:      1,
			// 		PeriodSeconds:       20,
			// 		SuccessThreshold:    1,
			// 		FailureThreshold:    5,
			// 		InitialDelaySeconds: 20,
			// 	},
		},
	}

	return containers, nil
}

func expectedShutdownManagerImage() string {
	if v := version.Get().ShutdownManagerVersion; v != "" {
		return fmt.Sprintf("%s:%s", strings.Split(egv1a1.DefaultShutdownManagerImage, ":")[0], v)
	}
	return egv1a1.DefaultShutdownManagerImage
}

func expectedShutdownManagerArgs(cfg *egv1a1.ShutdownConfig) []string {
	args := []string{"envoy", "shutdown-manager"}
	if cfg != nil && cfg.DrainTimeout != nil {
		args = append(args, fmt.Sprintf("--ready-timeout=%.0fs", cfg.DrainTimeout.Seconds()+10))
	}
	return args
}

func expectedShutdownPreStopCommand(cfg *egv1a1.ShutdownConfig) []string {
	command := []string{"envoy-gateway", "envoy", "shutdown"}

	if cfg == nil {
		return command
	}

	if cfg.DrainTimeout != nil {
		command = append(command, fmt.Sprintf("--drain-timeout=%.0fs", cfg.DrainTimeout.Seconds()))
	}

	if cfg.MinDrainDuration != nil {
		command = append(command, fmt.Sprintf("--min-drain-duration=%.0fs", cfg.MinDrainDuration.Seconds()))
	}

	return command
}

// expectedContainerVolumeMounts returns expected proxy container volume mounts.
func expectedContainerVolumeMounts(containerSpec *egv1a1.KubernetesContainerSpec) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		// {
		// 	Name:      "certs",
		// 	MountPath: "/certs",
		// 	ReadOnly:  true,
		// },
		// {
		// 	Name:      "sds",
		// 	MountPath: "/sds",
		// },
		{
			Name:      "ratelimiter-truststore-secret-volume",
			MountPath: "/home/wso2/security/truststore/ratelimiter-ca.crt",
			SubPath:   "ca.crt",
		},
		{
			Name:      "ratelimiter-truststore-secret-volume",
			MountPath: "/home/wso2/security/truststore/ratelimiter.crt",
			SubPath:   "tls.crt",
		},
		{
			Name:      "log-conf-volume",
			MountPath: "/home/wso2/conf/",
		},
		{
			Name:      "enforcer-keystore-secret-volume",
			MountPath: "/home/wso2/security/truststore/enforcer.crt",
			SubPath:   "tls.crt",
		},
		{
			Name:      "adapter-truststore-secret-volume",
			MountPath: "/home/wso2/security/truststore/adapter.crt",
			SubPath:   "tls.crt",
		},
		{
			MountPath: "/home/wso2/security/keystore/router.crt",
			Name:      "router-keystore-secret-volume",
			SubPath:   "tls.crt",
		},
		{
			MountPath: "/home/wso2/security/keystore/router.key",
			Name:      "router-keystore-secret-volume",
			SubPath:   "tls.key",
		},
	}

	return resource.ExpectedContainerVolumeMounts(containerSpec, volumeMounts)
}

func expectedEnforcerPorts() []corev1.ContainerPort {
	ports := []corev1.ContainerPort{
		{
			Name:          "ext-auth",
			ContainerPort: int32(8081),
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          "jwks",
			ContainerPort: int32(9092),
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          "cc-xds",
			ContainerPort: int32(18002),
			Protocol:      corev1.ProtocolTCP,
		},
	}

	for _, port := range config.ReadConfigs().Deployment.Gateway.EnforcerPorts {
		ports = append(ports,
			corev1.ContainerPort{
				Name:          port.Name,
				ContainerPort: port.ContainerPort,
				Protocol:      corev1.ProtocolTCP,
			},
		)
	}

	return ports
}

// expectedContainerVolumeMounts returns expected proxy container volume mounts.
func expectedEnforcerVolumeMounts(containerSpec *egv1a1.KubernetesContainerSpec) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "tmp",
			MountPath: "/tmp",
			ReadOnly:  true,
		},
		{
			Name:      "enforcer-keystore-secret-volume",
			MountPath: "/home/wso2/security/keystore/enforcer.key",
			SubPath:   "tls.key",
		},

		{MountPath: "/home/wso2/security/keystore/enforcer.crt",
			Name:    "enforcer-keystore-secret-volume",
			SubPath: "tls.crt",
		},
		{MountPath: "/home/wso2/security/truststore/apk.crt",
			Name:    "enforcer-keystore-secret-volume",
			SubPath: "ca.crt",
		},
		{MountPath: "/home/wso2/security/truststore/enforcer.crt",
			Name:    "enforcer-keystore-secret-volume",
			SubPath: "tls.crt",
		},
		{MountPath: "/home/wso2/security/truststore/adapter.crt",
			Name:    "adapter-truststore-secret-volume",
			SubPath: "tls.crt",
		},
		{MountPath: "/home/wso2/security/truststore/router.crt",
			Name:    "router-keystore-secret-volume",
			SubPath: "tls.crt",
		},
		{MountPath: "/home/wso2/conf/",
			Name: "log-conf-volume",
		},
		{MountPath: "/home/wso2/security/keystore/mg.pem",
			Name:    "enforcer-jwt-secret-volume",
			SubPath: "mg.pem",
		},
		{MountPath: "/home/wso2/security/truststore/mg.pem",
			Name:    "enforcer-jwt-secret-volume",
			SubPath: "mg.pem",
		},
		{MountPath: "/home/wso2/security/keystore/mg.key",
			Name:    "enforcer-jwt-secret-volume",
			SubPath: "mg.key",
		},
		{MountPath: "/home/wso2/security/truststore/wso2carbon.pem",
			Name:    "enforcer-trusted-certs",
			SubPath: "wso2carbon.pem",
		},
		{MountPath: "/home/wso2/security/truststore/wso2-apim-carbon.pem",
			Name:    "enforcer-apikey-cert",
			SubPath: "wso2-apim-carbon.pem",
		},
		{MountPath: "/home/wso2/security/truststore/idp.pem",
			Name:    "idp-certificate-secret-volume",
			SubPath: "wso2carbon.pem",
		},
	}

	return resource.ExpectedContainerVolumeMounts(containerSpec, volumeMounts)
}

// expectedDeploymentVolumes returns expected proxy deployment volumes.
func expectedDeploymentVolumes(deploymentSpec *egv1a1.KubernetesDeploymentSpec) []corev1.Volume {
	conf := config.ReadConfigs()
	volumes := []corev1.Volume{
		{
			Name: "ratelimiter-truststore-secret-volume",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  conf.Deployment.Gateway.Volumes.RatelimiterTruststoreSecretVolume,
					DefaultMode: ptr.To[int32](420),
				},
			},
		},
		{
			Name: "enforcer-keystore-secret-volume",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: ptr.To[int32](420),
					SecretName:  conf.Deployment.Gateway.Volumes.EnforcerKeystoreSecretVolume,
				},
			},
		},
		{
			Name: "log-conf-volume",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: ptr.To[int32](420),
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "apk-test-wso2-apk-log-conf",
					},
				},
			},
		},
		{
			Name: "router-keystore-secret-volume",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: ptr.To[int32](420),
					SecretName:  conf.Deployment.Gateway.Volumes.RouterKeystoreSecretVolume,
				},
			},
		},
		{
			Name: "adapter-truststore-secret-volume",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: ptr.To[int32](420),
					SecretName:  conf.Deployment.Gateway.Volumes.AdapterTruststoreSecretVolume,
				},
			},
		},
		{
			Name: "enforcer-jwt-secret-volume",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: ptr.To[int32](420),
					SecretName:  conf.Deployment.Gateway.Volumes.EnforcerJwtSecretVolume,
				},
			},
		},
		{
			Name: "enforcer-trusted-certs",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: ptr.To[int32](420),
					SecretName:  conf.Deployment.Gateway.Volumes.EnforcerTrustedCerts,
				},
			},
		},
		{
			Name: "enforcer-apikey-cert",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: ptr.To[int32](420),
					SecretName:  conf.Deployment.Gateway.Volumes.EnforcerApikeyCert,
				},
			},
		},
		{
			Name: "idp-certificate-secret-volume",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: ptr.To[int32](420),
					SecretName:  conf.Deployment.Gateway.Volumes.IDPCertificateSecretVolume,
				},
			},
		},
		{
			Name: "tmp",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}

	return resource.ExpectedDeploymentVolumes(deploymentSpec.Pod, volumes)
}

// expectedContainerEnv returns expected proxy container envs.
func expectedContainerEnv(containerSpec *egv1a1.KubernetesContainerSpec) []corev1.EnvVar {
	env := []corev1.EnvVar{
		{
			Name: envoyNsEnvVar,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "metadata.namespace",
				},
			},
		},
		{
			Name: envoyPodEnvVar,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "metadata.name",
				},
			},
		},
	}

	if containerSpec != nil {
		return resource.ExpectedContainerEnv(containerSpec, env)
	}
	return env

}

func expectedEnforcerEnv(containerSpec *egv1a1.KubernetesContainerSpec) []corev1.EnvVar {
	conf := config.ReadConfigs()
	env := []corev1.EnvVar{
		{
			Name:  "ADAPTER_HOST_NAME",
			Value: conf.Deployment.Gateway.AdapterHostName,
		},
		{
			Name:  "ADAPTER_HOST",
			Value: conf.Deployment.Gateway.AdapterHost,
		},
		{
			Name:  "COMMON_CONTROLLER_HOST_NAME",
			Value: conf.Deployment.Gateway.CommonControllerHostName,
		},
		{
			Name:  "COMMON_CONTROLLER_HOST",
			Value: conf.Deployment.Gateway.CommonControllerHost,
		},
		{
			Name:  "ENFORCER_PRIVATE_KEY_PATH",
			Value: conf.Deployment.Gateway.EnforcerPrivateKeyPath,
		},
		{
			Name:  "ENFORCER_PUBLIC_CERT_PATH",
			Value: conf.Deployment.Gateway.EnforcerPublicCertPath,
		},
		{
			Name:  "ENFORCER_SERVER_NAME",
			Value: conf.Deployment.Gateway.EnforcerServerName,
		},
		{
			Name:  "TRUSTED_CA_CERTS_PATH",
			Value: conf.Deployment.Gateway.AdapterTrustedCAPath,
		},
		{
			Name:  "ADAPTER_XDS_PORT",
			Value: "18000",
		},
		{
			Name:  "COMMON_CONTROLLER_XDS_PORT",
			Value: conf.Deployment.Gateway.CommonControllerXDSPort,
		},
		{
			Name:  "COMMON_CONTROLLER_REST_PORT",
			Value: conf.Deployment.Gateway.CommonControllerRestPort,
		},
		{
			Name:  "ENFORCER_LABEL",
			Value: conf.Deployment.Gateway.EnforcerLabel,
		},
		{
			Name:  "ENFORCER_REGION",
			Value: conf.Deployment.Gateway.EnforcerRegion,
		},
		{
			Name:  "XDS_MAX_MSG_SIZE",
			Value: conf.Deployment.Gateway.EnforcerXDSMaxMsgSize,
		},
		{
			Name:  "XDS_MAX_RETRIES",
			Value: conf.Deployment.Gateway.EnforcerXDSMaxRetries,
		},
		{
			Name:  "JAVA_OPTS",
			Value: "-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5006 -Dhttpclient.hostnameVerifier=AllowAll -Xms512m -Xmx512m -XX:MaxRAMFraction=2",
		},
	}

	if containerSpec != nil {
		return resource.ExpectedContainerEnv(containerSpec, env)
	}
	return env

}
