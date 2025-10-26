#!/bin/bash
kubectl exec -it $(kubectl get pod -n devops -l app=mysql -o name) -n devops -- sh
