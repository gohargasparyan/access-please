package accessplease

import (
	"fmt"
	"github.com/gohargasparyan/access-please/common"
	apOkta "github.com/gohargasparyan/access-please/okta"
	"github.com/gohargasparyan/access-please/rolesandbindings"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strings"
)

type AccessPlease struct {
	//todo get context from kube config
	context    string
	apcache    *cache.Cache
	kubeClient *kubernetes.Clientset
}

func New(context string) (*AccessPlease, error) {
	apcache := cache.New(cache.NoExpiration, cache.NoExpiration)

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	//todo add option for incluster config too
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	common.Panic(err)

	kubeClient, err := kubernetes.NewForConfig(config)
	common.Panic(err)

	return &AccessPlease{
		apcache:    apcache,
		context:    context,
		kubeClient: kubeClient,
	}, nil
}

func (instance *AccessPlease) Run() error {
	contextName := instance.context
	kubeClient := instance.kubeClient

	oktaConfig := okta.NewConfig()
	oktaClient := okta.NewClient(oktaConfig, nil, nil)

	rolesandbindings.InitResourcesCache(instance.apcache)
	rolesandbindings.AddReadOnlyAccess(kubeClient, *instance.apcache)
	apOkta.AddGroup(oktaClient, apOkta.ReadOnlyOktaGroup)

	instance.watchNamespace(kubeClient, contextName, oktaClient, *instance.apcache)

	return nil
}

func (instance *AccessPlease) watchNamespace(kubeClient kubernetes.Interface, contextName string, oktaClient *okta.Client, apcache cache.Cache) {
	log.WithField("context", contextName).
		Info("context selected")

	watcher, err := kubeClient.CoreV1().Namespaces().Watch(metav1.ListOptions{})
	common.Panic(err)

	ch := watcher.ResultChan()
	for event := range ch {
		if namespace, ok := event.Object.(*v1.Namespace); ok {
			var namespaceName = namespace.Name
			var envPrefix = strings.Split(contextName, ".")[0]
			var groupName = fmt.Sprintf("%s-%s", envPrefix, namespaceName)

			switch event.Type {
			case watch.Added:
				log.Info("------------------" + groupName)
				//rolesandbindings.AddReadWriteAccess(kubeClient, namespaceName, groupName, apcache)
				//apOkta.AddGroup(oktaClient, groupName)
			case watch.Deleted:
				//apOkta.DeleteGroup(oktaClient, groupName)
			}
		}
	}
}
