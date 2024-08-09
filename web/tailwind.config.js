/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{html,js,vue}"],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: true,
    darkTheme: "dark"
  },
  darkMode: ['class', '[data-theme="dark"]']
}

