const fs = require('fs');

(async () => {
  // todo SAP HANA 2021 の DB が立ち上げた後に API を叩いた直後はなぜかうまくいかない
  const axios = require('axios');

  const env = JSON.parse(fs.readFileSync('../../env.json', 'utf8'))['sap'];
  console.log(env);

  // sap-client の指定がないとエラーになってしまう
  // sap-client=number が number（integer） を返すようになったらうまくいくようになった
  const params = {
    'sap-client': env.client,
  };

  const apiService = {
    odata: {
      v4: {
        bank: {
          detail: {
            url: 'sap/opu/odata4/sap/api_bank/srvd_a2x/sap/api_bank_2/0001/Bank',
            method: 'get',
            params: {
              $filter: `BankCountry eq '${env.bank.country}' and Bank eq '${env.bank.bankId}'`,
            }
          },
          list: {
            url: 'sap/opu/odata4/sap/api_bank/srvd_a2x/sap/api_bank_2/0001/Bank',
            method: 'get',
            params: {}
          }
        }
      },
      v2: {
        product: {
          detail: {
            url: 'sap/opu/odata/sap/API_PRODUCT_SRV/A_Product',
            method: 'get',
            params: {}
          }
        }
      }
    },
  };

  const version = 'v4';
  const service = 'bank';
  const type = 'detail';

  const createUrl = () => {
    return `http://${env.host}:${env.port}/${apiService.odata[version][service][type].url}`
  }

  const url = createUrl();

  const username = env.username;
  const password = env.password;

  const auth = {
    username,
    password,
  }

  // cookie を使用するように設定
  axios.defaults.withCredentials = true

  const xsrfToken = '';

  const headerConfig = {
    'Content-Type': 'application/json;odata.metadata=minimal;charset=utf-8',
    'Accept': 'application/json',
    'X-CSRF-Token': xsrfToken || 'Fetch',
    // user agent が必要かと思ったが実際は sap-client の指定があればうまくいった
    // 'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36',
  }

  const createAxiosConfig = () => {
    return {
      params: {
        ...params,
        ...apiService.odata[version][service][type].params,
      },
      auth,
      headers: {
        ...headerConfig,
      },
    }
  };

  let createAxiosConfigResult;

  try {
    createAxiosConfigResult = createAxiosConfig();

    console.log(`url: `, url);
    console.log(`createAxiosConfigResult: `, createAxiosConfigResult);
  } catch (e) {
    console.log(`creating axios config error`);
    console.log(e);
  }

  try {
    const response = await axios.get(
      url,
      createAxiosConfigResult,
      {}
    );

    console.log(`request is success`);
    console.log('response.data: ', response.data);
  } catch (error) {
    console.log(`request error message`)
    // set-cookie の値を確認する
    console.error(error.response.data)
  }
})();
