SUBDIRS= \
	src

.PHONY: prune-replicasets

default:
	(cd src; make)

prune-replicasets:
	@kubectl get replicasets | awk '$$2==0 { system("kubectl delete replicaset/" $$1) }'

test:
	@for d in ${SUBDIRS}; do \
		if [ -f $$d/Makefile ] ; then \
			if grep '^test:' $$d/Makefile ; then \
				echo $$d ; \
				(cd $$d ; make test) ; \
			fi \
		fi \
	done
