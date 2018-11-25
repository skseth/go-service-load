import http from "k6/http";
import { check, group, sleep } from "k6";
import { Rate } from "k6/metrics";

var failureRate = new Rate("check_failure_rate");

export let options = {
    vus: 5,
    collectors: {
        "influxdb": {
          "tagsAsFields": ["vu","iter", "url", "name"]
        }
      },
    stages: [
        { duration: "10s", target: 100 },
        { duration: "10s", target: 100 },
        { duration: "10s", target: 350 },
        { duration: "10s", target: 350 },
        { duration: "10s", target: 1000 },
        { duration: "50s", target: 1000 },
        { duration: "20s", target: 5000 },
        { duration: "50s", target: 5000 },
        { duration: "2m", target: 0 },
    ],
    thresholds: {
        // We want the 95th percentile of all HTTP request durations to be less than 500ms
        "http_req_duration": ["p(95)<500"],
        "check_failure_rate": [
            // Global failure rate should be less than 1%
            "rate<0.01",
            // Abort the test early if it climbs over 5%
            { threshold: "rate<=0.05", abortOnFail: true },
        ],
    },
};

export default function() {
    const response = http.get("http://localhost:8000/chunk/1");

    // check() returns false if any of the specified conditions fail
    let checkRes = check(response, {
        "status is 200": (r) => r.status === 200,
        "content is present": (r) => r.status === 200 && r.body.length > 10,
    });

    failureRate.add(!checkRes)

    sleep(1)
};