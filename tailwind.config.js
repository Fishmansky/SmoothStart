/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
	'./views/**/*.templ',
	'./views/**/*.go',
  ],
  theme: {
    screens: {
      sm: '480px',
      md: '768px',
      lg: '976px',
      xl: '1440px',
    },
    colors: {
	    'dark': '#0b2027',
	    'green': '#b5f44a',
	    'mint': '#70a9a1',
	    'blue': '#40798c',
	    'light': '#cfd7c7',
	    'gray': '#5e7572',
	    'red': '#d72638',
	    'cream': '#f6f1d1',
    },
    fontFamily: {
	    sans: ['Nunito Sans', 'sans-serif'],
    },
    extend: {},
  },
  plugins: [],
}

