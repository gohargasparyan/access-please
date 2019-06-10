# Access Please is a golang project to help you manage roles/bindings to access to kuberenetes

The idea is to dynamically generate roles and bindings in kuberenetes for each namespace for the given context. 
This tool can be used in companies, where teams create namespaces for their projects and need role based access to clusters. 
Here I use groups in {context}-{namespace} format , e.g. dev-namespaceA. There are 2 flavours of access, read-only and read-write. 
By default read-only cluster role is added. On namespace creation event read-write role will be added and bound to corresponding group. 

//todo make this optional
On top of that the very same group will be added to Okta, which is the chosen access managemen tool.

# Preconditions to Run
You need to have installed
* golang sdk
* install kubectl
* configure $HOME/.kube/config for access of desired cluster
* configure .okta/okta.yaml for okta api access

# Run
To run app run following command passing the context:

        $ go run . --context=context

