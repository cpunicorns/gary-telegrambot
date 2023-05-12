gcloud functions deploy gary-telegrambot \
	--gen2 \
	--region=europe-west3 \
	--runtime=go118 \
	--source=https://source.developers.google.com/p/famous-badge-386309/r/garybot \
	--entry-point=main \
	--env-vars-file=.env.yaml \