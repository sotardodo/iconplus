#!/bin/bash
kubectl exec -it $(kubectl get pod -n devops -l app=laravel -o name) -n devops -- sh
