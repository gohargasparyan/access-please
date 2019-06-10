package rolesandbindings

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr"
	"github.com/gohargasparyan/access-please/common"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	ResourcesDir                  = "../resources"
	ReadOnlyClusterRole           = "ReadOnlyClusterRole"
	ReadOnlyClusterRoleDir        = "roles/read-only-cluster-role.yml"
	ReadOnlyClusterRoleBinding    = "ReadOnlyClusterRoleBinding"
	ReadOnlyClusterRoleBindingDir = "bindings/read-only-binding.yml"
	ReadWriteRole                 = "ReadWriteRole"
	ReadWriteRoleDir              = "roles/read-write-role.yml"
	ReadWriteRoleBinding          = "ReadWriteRoleBinding"
)

func InitResourcesCache(apcache *cache.Cache) {
	box := packr.NewBox(ResourcesDir)

	// add read only cluster role to cache
	readOnlyClusterRoleBytes, err := box.Find(ReadOnlyClusterRoleDir)
	common.Panic(err)
	var clusterRole rbacV1.ClusterRole
	err = yaml.Unmarshal(readOnlyClusterRoleBytes, &clusterRole)
	common.Panic(err)
	apcache.Set(ReadOnlyClusterRole, clusterRole, cache.NoExpiration)

	// add read only cluster role binding to cache
	readOnlyClusterRoleBindingsBytes, err := box.Find(ReadOnlyClusterRoleBindingDir)
	common.Panic(err)
	var clusterRoleBinding rbacV1.ClusterRoleBinding
	err = yaml.Unmarshal(readOnlyClusterRoleBindingsBytes, &clusterRoleBinding)
	common.Panic(err)
	apcache.Set(ReadOnlyClusterRoleBinding, clusterRoleBinding, cache.NoExpiration)

	//add read/write role to cache
	readWriteRoleDirBytes, err := box.Find(ReadWriteRoleDir)
	common.Panic(err)
	var role rbacV1.Role
	err = yaml.Unmarshal(readWriteRoleDirBytes, &role)
	common.Panic(err)
	apcache.Set(ReadWriteRole, role, cache.NoExpiration)

	//ad read/write RoleBinding to cache
	roleBinding := rbacV1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		RoleRef: rbacV1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
		},
	}

	apcache.Set(ReadWriteRoleBinding, roleBinding, cache.NoExpiration)
}

//todo think of error handling
func AddReadOnlyAccess(client kubernetes.Interface, apcache cache.Cache) {
	addReadOnlyClusterRole(client, apcache)
	addReadOnlyClusterRoleBinding(client, apcache)
	log.Printf("Added read-only access to everyone in group 'global-read-only'.")
}

func AddReadWriteAccess(client kubernetes.Interface, namespace string, groupName string, apcache cache.Cache) {
	var roleName = fmt.Sprintf("%s-read-write-role", namespace)

	addReadWriteRole(client, namespace, roleName, apcache)
	addReadWriteRoleBinding(client, namespace, roleName, groupName, apcache)

	//todo use this log format everywhere
	log.WithField("groupName", groupName).
		Info("Added read-write access to everyone in group")
}

func addReadOnlyClusterRole(client kubernetes.Interface, apcache cache.Cache) {
	clusterRoleI, found := apcache.Get(ReadOnlyClusterRole)
	var clusterRole = clusterRoleI.(rbacV1.ClusterRole)

	if found {
		client.RbacV1().ClusterRoles().Update(&clusterRole)
		log.Printf("Added read-only ClusterRole.")
	} else {
		log.Printf("Read Only ClusterRole not found in cache. Can not add Read Only ClusterRole")
	}
}

func addReadOnlyClusterRoleBinding(client kubernetes.Interface, apcache cache.Cache) {
	clusterRoleBindingI, found := apcache.Get(ReadOnlyClusterRoleBinding)
	var clusterRoleBinding = clusterRoleBindingI.(rbacV1.ClusterRoleBinding)

	if found {
		client.RbacV1().ClusterRoleBindings().Update(&clusterRoleBinding)
		log.Printf("Added read-only ClusterRoleBinding.")
	} else {
		log.Printf("Read Only ClusterRoleBinding not found in cache. Can not add Read Only ClusterRoleBinding")
	}
}

func addReadWriteRole(client kubernetes.Interface, namespace string, roleName string, apcache cache.Cache) {
	roleI, found := apcache.Get(ReadWriteRole)
	var role = roleI.(rbacV1.Role)

	if found {
		var roleCopy = *role.DeepCopy()
		roleCopy.ObjectMeta.Namespace = namespace
		roleCopy.ObjectMeta.Name = roleName
		client.RbacV1().Roles(namespace).Update(&roleCopy)
		log.Printf("Added read-write Role for namespace : %v", namespace)
	} else {
		log.Printf("Read/Write Role not found in cache. Can not add Role for Namespace: %v", namespace)
	}
}

func addReadWriteRoleBinding(client kubernetes.Interface, namespace string, roleName string, groupName string, apcache cache.Cache) {
	roleBindingI, found := apcache.Get(ReadWriteRoleBinding)
	var roleBinding = roleBindingI.(rbacV1.RoleBinding)

	if found {
		var roleBindingCopy = *roleBinding.DeepCopy()
		roleBindingCopy.ObjectMeta = metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-read-write-role-binding", namespace),
			Namespace: namespace,
		}
		roleBindingCopy.Subjects = []rbacV1.Subject{{
			Kind: "Group",
			Name: groupName,
		}}
		roleBindingCopy.RoleRef.Name = roleName

		client.RbacV1().RoleBindings(namespace).Update(&roleBinding)
	} else {
		log.Printf("Read/Write RoleBinding not found in cache. Can not add RoleBinding for Namespace: %v", namespace)
	}
}
