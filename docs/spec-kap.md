# 📜 KubeAegisPolicy Specification

```yaml
apiVersion: cclab.kubeaegis.com/v1
kind: KubeAegisPolicy
metadata:
  name: [policy name]
  namespace: [namespace name]
spec:
  enableReport: [true|false]
  requestRule:
    - type: [network|system|cluster]
      selector:
        kind: [pod|namespace                 
              |service|deployment]
        namespace: [namespace name]
        matchLabels:
        [key1]: [value1]
        cel:
         - [cel expression]
      rule:
        action: [Allow|Block|Log|              
                |Trace|Enforce|Audit]
        from:
          - kind: [endpoint|entities|namespace|
                  |serviceAccounts|cidr|port
                  |protocol|fqdns]
            labels:
              - [key1]: [value1]
            args: [<arg1>, <arg2>, ...]
            port: [port number]
            protocol: [TCP|UDP|ICMP]
        to:
          - kind: [endpoint|namespace|
                  |serviceAccounts|entities
                  |cidr|port|protocol|fqdns]
            labels:
              - [key1]: [value1]
            args: [<arg1>, <arg2>, ...]
            port: [port number]
            protocol: [TCP|UDP|ICMP]
        actionPoint:
          - subType: [http|file|process|network|
                     |syscalls|capabilities|
                     |generate|mutate|validate|
                     |verifyImage|cleanup]
            resource:           
                path: [resource path]
                pattern: [pattern]
                kind: [resource kind]
                filter:
                  - condition: [any|all]
                    key: [filter key]
                    operator: [Equals|In
                              |NotIn|Exists]
                    value: [filter value]
                details:
                  - [key1]: [value1]
                  - [key2]: [value2]
status:
  status: [policy enforcement status]
  lastUpdated: [last update time] 
  numberOfAPs: [number]
  listOfAPs: [<name1>, <name2>, ...]
  numberOfTargets: [number]
  listOfTargets: [<name1>, <name2>, ...]
```