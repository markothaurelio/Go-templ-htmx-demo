/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.templ", "./static/**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [require('@tailwindcss/typography')],
}

