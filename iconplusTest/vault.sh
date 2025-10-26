#!/bin/bash
kubectl exec -it $(kubectl get pod -n devops -l app=vault -o name) -n devops -- sh
