export const BASE_URL = "http://localhost";
export const url = (path) => `${BASE_URL}${path}`;
export const getRandomAccount = () => {
  const accounts = ["mary", "sandra", "christina", "pauline"];
  const INDEX = Math.floor(Math.random() * accounts.length);
  return [accounts[INDEX], accounts[INDEX] + accounts[INDEX]];
};
