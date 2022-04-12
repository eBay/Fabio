---
title: "registry.consul.namespace"
---

`registry.consul.namespace` configures the consul namespace in which fabio will register itself.

 Namespaces are a feature only available in the enterprise version of consul. In the open-source
 version or with an empty namespace option fabio will be registered in the default namespace. Only
 services running in the same consul namespace will be picked up by fabio.

The default is

	registry.consul.namespace = 
