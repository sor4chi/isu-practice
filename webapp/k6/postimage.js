import http from "k6/http";
import { getRandomAccount, url } from "./config.js";
import { parseHTML } from "k6/html";

const testImage = open("./test.png", "b");

export default function () {
  const [ACCOUNT_NAME, PASSWORD] = getRandomAccount();
  const loginRes = http.post(url("/login"), {
    account_name: ACCOUNT_NAME,
    password: PASSWORD,
  });

  const doc = parseHTML(loginRes.body);

  const token = doc.find("input[name=csrf_token]").first().attr("value");

  http.post(url("/"), {
    csrf_token: token,
    file: http.file(testImage, "test.png", "image/png"),
    body: "JAVA SCRIPT API",
  });
}
