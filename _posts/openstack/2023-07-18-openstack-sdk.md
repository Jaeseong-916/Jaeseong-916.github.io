---
title: "OpenStack SDK 코드 분석 - server list 호출 과정"
date: 2023-07-18
categories:
  - openstack
tags:
  - openstack
  - openstacksdk
  - python
  - 코드분석
---

OpenStack CLI에서 `server list`를 호출할 때, 내부적으로 OpenStack SDK가 어떻게 동작하는지 분석한 내용입니다.

## 시작점: ListServer 클래스

```python
class ListServer(command.Lister):
    _description = _("List servers")

    def take_action(self, parsed_args):
        compute_client = self.app.client_manager.sdk_connection.compute
        data = list(compute_client.servers(**search_opts))
        table = (
            column_headers,
            (
                utils.get_item_properties(...) for s in data
            ),
        )
        return table
```

핵심은 아래 두 줄입니다:

```python
compute_client = self.app.client_manager.sdk_connection.compute
data = list(compute_client.servers(**search_opts))
```

## 따라가 봅시다!

### compute_client.servers

이 메서드는 `openstacksdk.openstack.compute.v2._proxy.Proxy.servers`를 호출합니다.

```python
def servers(self, details=True, all_projects=False, **query):
    if all_projects:
        query['all_projects'] = True
    base_path = '/servers/detail' if details else None
    return self._list(_server.Server, base_path=base_path, **query)
```

### _list

`openstacksdk.openstack.proxy.Proxy._list`에서 동작합니다.

```python
def _list(self, resource_type, paginated=True, base_path=None,
          jmespath_filters=None, **attrs):
    data = resource_type.list(
        self, paginated=paginated, base_path=base_path, **attrs
    )
    if jmespath_filters and isinstance(jmespath_filters, str):
        return jmespath.search(jmespath_filters, data)
    return data
```

### _server.Server

`openstacksdk.openstack.compute.v2.server.Server` 클래스를 호출하며, `resource.Resource`를 이용한 변수를 재정의하고 있습니다.

### Resource

`openstacksdk.openstack.resource.Resource` 클래스는 API request/response에 사용되는 Body를 구성하고, API CRUD의 주체가 되는 함수들이 정의되어 있습니다.

## 종합

| 클래스 | 역할 |
|--------|------|
| **_proxy** | 상위 수준의 CRUD 함수가 정의, 파라미터를 통해 proxy 호출 |
| **proxy** | 매개변수로 전달받은 파라미터를 이용하여 리소스를 가공하여 제공 |
| **Server** | 상속한 Resource 클래스의 변수를 재정의 |
| **Resource** | CRUD 과정의 함수 정의, API request/response에 필요한 리소스 관리 |
