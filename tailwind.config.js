/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./http/templates/*.{html,js}"],
    theme: {
      extend: {},
    },
    plugins: [
      require('@tailwindcss/forms'),
      require('@tailwindcss/typography'),
    ],
}