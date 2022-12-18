## What is cloudctl?
This project aims to take control of different cloud resources with kubernetes operator. By 
defining the CRD(CustomResourceDefinition) of the cloud resource, it's easy to execute the underlying
cloud providers' API and get the status of the cloud resource. 

## Usage
### define CRD
1.  _lifecycle_ is an object type and is used to save the cloud API json file.
2.  _domain_ is used to save underlying cloud resource metadata. 
3. Users need to add some metadata of the cloud resouce, like UUID, which is used to find the specfic resource. 
4. Users need to define a secret which contains the authentication information. And the _secretRef_ contains the secret's name and namespace.
```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: openstackservers.doslab.io
spec:
  group: doslab.io
  names:
    kind: OpenstackServer
    plural: openstackservers
    singular: openstackserver
  scope: Namespaced
  versions:
    - additionalPrinterColumns:
        - description: Server ID
          jsonPath: .spec.id
          name: ID
          type: string
        - description: Server Name
          jsonPath: .spec.domain.name
          name: Name
          type: string
        - description: Server Status
          jsonPath: .spec.domain.status
          name: Status
          type: string
        - description: host id is where the server is located in the cloud
          jsonPath: .spec.domain.hostid
          name: HostID
          type: string
        - description: image
          jsonPath: .spec.domain.image.id
          name: ImageID
          type: string
        - description: flavor
          jsonPath: .spec.domain.flavor.id
          name: FlavorID
          type: string
      name: v1
      schema:
        openAPIV3Schema:
          description: VMInstance is the Schema for the vminstances API
          properties:
            apiVersion:
              type: string
            kind:
              type: string
            metadata:
              type: object
            spec:
              description: spec defines the desired state of openstack server
              properties:
                domain:
                  type: object
                  x-kubernetes-preserve-unknown-fields: true
                lifeCycle:
                  description: request to be execute
                  type: object
                  x-kubernetes-preserve-unknown-fields: true
                #                  metadata
                id:
                  type: string
                #                  secret info requeired
                secretRef:
                  description: SrereteRef
                  properties:
                    name:
                      description: secretName
                      type: string
                    namespace:
                      description: secretNamespace
                      type: string
                  required:
                    - name
                    - namespace
                  type: object
              required:
                - secretRef
              type: object
            status:
              type: object
          type: object
      served: true
      storage: true
      subresources:
        status: {}

```

### define CRD config 
```json
{
  "CrdName": "OpenstackServer",
  "MetaInfos": [
    {
      "SpecName": "id", // name in CRD spec
      "DomainName": "id", // name in CRD domain
      "CloudParaName":"id", // name in Cloud API parameter
      "InitJsonPath": "GetComputeV2Servers.id", 
      "DeleteJsonPath": "DeleteComputeV2Servers.id",
      "InitRespJsonPath": "server.id", //using to get the resource metadata from the GET Response
      "IsArray": false //whether the GET Response is array or only one resource
    }
  ],
  "InitJson": { //API name and essential parameters to get the sepcfic cloud resource
    "GetComputeV2Servers": {
      "id": ""
    }
  },
  "DeleteJson": { //API name and essential parameters to delete the sepcfic cloud resource
    "DeleteComputeV2Servers": {
      "id": ""
    }
  },
  "DomainJsonPath":"server" //used to get the domain info from the GET API response.
}
```

## deploy the cloudctl in your kubernetes cluster
```shell
make install
```

### add CRD to kubernets cluster
```shell
kubectl apply -f OpenstackServer.yaml
```

### Execute Cloud API
The json below is used to Create a Openstack Server. The response will contains the server's UUID and 
will be filled in the _spec.id_. The controller will automatically execute the GET Request from 
the config's InitJson, and fill the _domain_ with the response.
```json
{
  "apiVersion": "doslab.io/v1",
  "kind": "OpenstackServer",
  "metadata": {
    "name": "openstack-server-create"
  },
  "spec": {
    "lifeCycle": {
      "CreateComputeV2Servers": {
        "Opts": {
          "Name": "test-create",
          "ImageRef": "952b386b-6f30-46f6-b019-f522b157aa3a",
          "FlavorRef": "3"
        }
      }
    },
    "id": "",
    "secretRef": {
      "namespace": "default",
      "name": "openstack-compute-secret"
    }
  }
}
```

## Support APIS
### Openstack

| Resource Name | yaml | config| Create   | Get | Update | Delete | document |
|---------------|-----|-------|----------|---------|----------|--------|----------|
| Server        |     | Create | Get      | Update | Delete | document ||
| Image         |  | Create | Get      | Update | Delete | document ||
| Network       |  | Create | Get      | Update | Delete | document ||
| Snapshot      |  | Create | Get      | Update | Delete | document ||
| Route         |  | Create | Get      | Update | Delete | document ||
