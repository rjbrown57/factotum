# https://devhints.io/makefile

CONFIG=kwok-config.yaml
NODE_COUNT=50
NAME="kwok"

.PHONY: default
default: create

.PHONY: create
create:
	@echo "Creating kwok cluster"
	kwok create cluster -c $(CONFIG)
	kubectl apply -f ../config/crd/bases/

.PHONY: delete
delete:
	@echo "Deleting kwok cluster"
	kwok delete cluster --name $(NAME)

.PHONY: scale
scale: 
	@echo "Scaling kwok cluster"
	kwok scale node --replicas $(NODE_COUNT) --name $(NAME)

