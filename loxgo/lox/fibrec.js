const fib = (n) => {
  if (n <= 1) return n;
  return fib(n - 2) + fib(n - 1);
};

const start = new Date().getTime();
for (let i = 0; i < 30; i = i + 1) {
  console.log(fib(i));
}
const end = new Date().getTime();

console.log("TIME");
console.log((end - start) / 1000);
