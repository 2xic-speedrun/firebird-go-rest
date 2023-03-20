"""
One rest endpoint is running firebird-rest (nodejs)
One rest endpoint is running this implementation.

We want't to compare the latency
"""
import requests
import json 
import time
import asyncio
import matplotlib.pyplot as plt
import numpy as np
import math
import time
from scipy.stats import norm

SAMPLES = 1_0
REQUESTS = 10
LOCAL_ONLY = False

async def send_request(
    name,
    url,
    payload
):
    start = time.time()
    status = requests.post(
        url,
        json=payload
    )
    end = time.time()
    time_ms = (end - start) * 1_000
    
    return (
        status.status_code,
        time_ms
    )

def gaussian(time_usage, label):
    print(time_usage)
    mu = np.mean(time_usage)
    variance = np.var(time_usage)
    sigma = math.sqrt(variance)

    x = np.linspace(mu - 4*sigma, mu + 4*sigma, 100)
    plt.plot(x, norm.pdf(x, mu, sigma), label=label)

async def run_test(results, data, test_name, payloads):
    results[test_name] = {}
    for _ in range(SAMPLES):
        for _ in range(REQUESTS):
            calls = []
            names = []
            for name, rest_endpoints in data["configs"].items():
                if LOCAL_ONLY and not "local" in name.lower():
                    continue
                elif "local" in name.lower():
                    continue

                if name not in results[test_name]:
                    results[test_name][name] = []
                combined = rest_endpoints
                for key, value in payloads.items():
                    combined[key] = value

                calls.append(
                    (lambda host: (lambda:send_request(
                        name,
                        host,
                        combined
                    ))) (rest_endpoints["host"])
                )
                names.append(name)                
        values = await asyncio.gather(*[
            i() for i in calls
        ])

        for (name, i) in zip(names, values):
            (status_code, time_ms) = i
            results[test_name][name].append(time_ms)


async def main(): 
    with open("benchmark_config.json", "r") as file:
        data = json.load(file)
        results = {}
        for test_name, payloads in data["statements"].items():
            await run_test(results, data, test_name, payloads)
            plt.clf()
            print(test_name)
            for key, value in list(results[test_name].items()):
                print(value)
                gaussian(
                    value,
                    label=key
                )
            plt.legend()
            plt.savefig("benchmark/" + test_name.lower().replace(" ", "_") + ".png")

if __name__ == "__main__":
    asyncio.run(main())
