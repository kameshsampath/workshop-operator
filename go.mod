module github.com/kameshsampath/workshop-operator

require (
	github.com/eclipse/che-operator v0.0.0-20190905134656-e6153050c26d
	github.com/emicklei/go-restful v2.10.0+incompatible // indirect
	github.com/go-logr/logr v0.1.0
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/go-openapi/spec v0.19.3
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/openshift/api v3.9.1-0.20190916204813-cdbe64fb0c91+incompatible
	//plm release-4.3
	github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190914140755-ceebad14a2f8
	github.com/operator-framework/operator-marketplace v0.0.0-20190912222040-fff7b83b4e24
	github.com/operator-framework/operator-sdk v0.10.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20190911031432-227b76d455e7
	golang.org/x/net v0.0.0-20190918130420-a8b05e9114ab // indirect
	golang.org/x/tools v0.0.0-20190919031856-7460b8e10b7e // indirect
	k8s.io/api v0.0.0-20190919035539-41700d9d0c5b
	k8s.io/apimachinery v0.0.0-20190917163033-a891081239f5
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/gengo v0.0.0-20190907103519-ebc107f98eab // indirect
	k8s.io/klog v0.4.0 // indirect
	k8s.io/kube-openapi v0.0.0-20190918143330-0270cf2f1c1d
	sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools v0.1.10
)

// Pinned to kubernetes-1.13.4
replace (
	k8s.io/api => k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190228180357-d002e88f6236
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190228174230-b40b2a5939e4
)

// Pinned to latest OpenShift that can use OpenShift 4 API (release-4.3)
replace (
	github.com/openshift/api => github.com/openshift/api v3.9.1-0.20190916204813-cdbe64fb0c91+incompatible
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20190813201236-5a5508328169
)

replace (
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.29.0
	// Pinned to v2.9.2 (kubernetes-1.13.1) so https://proxy.golang.org can
	// resolve it correctly.
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.8.2-0.20190424153033-d3245f150225
	k8s.io/kube-state-metrics => k8s.io/kube-state-metrics v1.6.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v0.0.0-20190820212518-960c3cc04183
)

//replace github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.10.0

go 1.13
