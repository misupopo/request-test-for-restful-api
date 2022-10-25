host := $(shell jq -r .sap.host env.json)
port := $(shell jq -r .sap.port env.json)
username := $(shell jq -r .sap.username env.json)
password := $(shell jq -r .sap.password env.json)
client := $(shell jq -r .sap.client env.json)
bankCountry := $(shell jq -r .sap.bank.country env.json)
bankId := $(shell jq -r .sap.bank.bankId env.json)
xCSRFToken := $(shell jq -r .sap.xCSRFToken env.json)

# xCSRFToken に値がセットされていなければ Fetch がセットされる
ifeq ($(xCSRFToken), )
	xCSRFToken:=Fetch
endif

create-require-files:
	mkdir -p cookies
	mkdir -p headers
	touch cookies/cookies.txt
	touch headers/headers.txt

test: create-require-files
	echo "xCSRFToken: $(xCSRFToken)"
	echo "client: $(client)"
	echo "bankCountry: $(bankCountry)"
	echo "bankId: $(bankId)"
	echo "bankId: $$filter"

# 正常なレスポンスが返ってこない場合
# query parameter がついているかどうかを nginx の access.log で確認すること
get-bank-detail: create-require-files
	curl -s -i -c cookies/cookies.txt -b cookies.txt \
        -X GET -G \
        --data-urlencode "$$filter=BankCountry eq '$(bankCountry)' and Bank eq '$(bankId)'" \
        --data-urlencode 'sap-client=$(client)' \
        -H "Content-Type: application/json" \
        -H "Accept: application/json" \
        -H "X-CSRF-Token: $(xCSRFToken)" \
        -D headers/headers.txt \
        -u $(username):$(password) \
        "http://$(host):$(port)/sap/opu/odata4/sap/api_bank/srvd_a2x/sap/api_bank_2/0001/Bank"

get-bank-list: create-require-files
	curl -s -i -c cookies/cookies.txt -b cookies.txt \
        -X GET -G \
        --data-urlencode 'sap-client=$(client)' \
        -H "Content-Type: application/json" \
        -H "Accept: application/json" \
        -H "X-CSRF-Token: $(xCSRFToken)" \
        -D headers/headers.txt \
        -u $(username):$(password) \
        "http://$(host):$(port)/sap/opu/odata4/sap/api_bank/srvd_a2x/sap/api_bank_2/0001/Bank"
