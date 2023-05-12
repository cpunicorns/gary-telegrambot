#!/bin/bash

gcloud functions deploy garybot \
	--region=europe-west3 \
	--runtime=go118 \
	--source=https://source.developers.google.com/projects/famous-badge-386309/repos/garybot/fixed-aliases/ \
	--entry-point=main \
	--env-vars-file=.env.yaml \
	--trigger-http