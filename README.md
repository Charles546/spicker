Stock Price Grabber
===================

This is a simple API built with [go-swagger](https://goswagger.io). It makes a call to [alphavantage API](https://www.alphavantage.co/documentation/) to returns a JSON object with processed data.


Testing
-------

You can use command `go test -v ./...` or use the wrapper script `scripts/test.sh` to run the tests.

Running locally
---------------

A wrapper script `scripts/run.sh` is provided to easily launch the API server. You can defined the required environment variables in `.env` or `.default.env` files. Below are some of the environment variables you can use.

 - **SYMBOL**: The stock ticker symbol you want to query
 - **NDAYS**: The number of days you want to query
 - **ALPHAVANGAGE_APIKEY**: The API key to [alphavantage.co](https://alphavantage.co)
 - **REDIS_CONNECTION**: A url to a redis instance, such as `redis://@localhost:6379/0`, optional

You can leave everything except the APIKEY to default (See `.default.env`). You can supply it through the command, or add the value to your `.env` file. For example
```bash
ALPHAVANTAGE_APIKEY=xxxxxxx ./scripts/run.sh
```

Once the server is running, you can access the swagger documents at [http://localhost:3000/docs](http://localhost:3000/docs).

Run the below curl command to see the result.

```bash
curl http://localhost:3000/stockprices
```

Building
--------

The wrapper script `scripts/build.sh` will build the API server into a docker image and push into the docker registry. It requires **IMAGE_REPO** and **IMAGE_TAG** environment variables. For pushing the image, Docker hub credential needs to be configured.


Deploying
---------

The wrapper script `scripts/deploy.sh` will deploy the application into a Kubernetes cluster, including a deployment, a service, a secret and a horizontal pod autoscaler. You can customize the deployments with all the environment variables. Besides the environment variables mentioned above in the [Running locally](#running-locally) section, a few more environments are supported.

 - **ALPHAVANTAGE_APIKEY_BASE64**: base64 encoded API key used for creating the Kubernetes secret
 - **KUBE_CONTEXT**: Specify the kubectl context explicitly to avoid mistakes
 - **IMAGE_REPO**: The image repo used for the deployment
 - **IMAGE_TAG**: The tag of the image used for the deployment
 - **MIN_REPLCA**: Autoscaler allowed minimal number of replcas
 - **MAX_REPLCA**: Autoscaler allowed maximum number of replcas
 - **CPU_REQUEST**: The number of the CPU requested at the pod startup
 - **CPU_LIMIT**: The limit of how much CPU the pod can use
 - **CPU_THRESHOLD**: The threshold used for measuring the load for autoscaler
 - **MEM_REQUEST**: The size of the memory requested at the pod startup
 - **MEM_LIMIT**: The limit of how much memory the pod can use
 - **MEM_THRESHOLD**: The threshold used for measuring the memory utilization for autoscaler

The environment variables should go in the `.env.<environment>` file such as `.env.dev`, `.env.prod` etc. Specify the environment with `ENV` variable.

For example, supply the key and kubectl context through command line:
```bash
ALPHAVANTAGE_APIKEY_BASE64=xxxxxxxx KUBE_CONTEXT=`kubectl config current-context` ENV=dev ./scripts/deploy.sh
```

To take advantage of redis caching, assuming you have everything else in the `.env.prod` file:
```bash
kubectl apply -f auxiliary/redis.yaml
REDIS_CONNECTION=redis://@stock-prices-grabber-cache:6379/0  ENV=prod ./scripts/deploy.sh
```

