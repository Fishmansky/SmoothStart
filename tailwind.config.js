/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
	'./views/**/*.templ',
  ],
  theme: {
    screens: {
      sm: '480px',
      md: '768px',
      lg: '976px',
      xl: '1440px',
    },
    colors: {
	    'mywhite': '#ebf2fa',
	    'dark': '#25283d',
	    'green': '#297373',
	    'yellow': '#bdd358',
	    'violet': '#564256',
    },
    fontFamily: {
	    sans: ['Nunito Sans', 'sans-serif'],
    },
    extend: {},
  },
  plugins: [],
}

