import http from "k6/http";
import { BASE_URL } from "./config.js";

export default function () {
  http.get(`${BASE_URL}/`);
}
