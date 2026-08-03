/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts}'],
  theme: {
    extend: {
      colors: {
        ef: {
          // Backgrounds (dark hard variant)
          'bg-dim':    '#232a2e',
          'bg0':       '#2b3339',
          'bg1':       '#323d43',
          'bg2':       '#3a464c',
          'bg3':       '#434f55',
          'bg4':       '#4f5b62',
          // Tinted backgrounds
          'bg-red':    '#514045',
          'bg-green':  '#425047',
          'bg-yellow': '#4d4c43',
          'bg-blue':   '#3a515d',
          // Foreground
          'fg':        '#d3c6aa',
          'grey0':     '#7a8478',
          'grey1':     '#859289',
          'grey2':     '#9da9a0',
          // Accents
          'red':       '#e67e80',
          'orange':    '#e69875',
          'yellow':    '#dbbc7f',
          'green':     '#a7c080',
          'aqua':      '#83c092',
          'blue':      '#7fbbb3',
          'purple':    '#d699b6',
        },
      },
    },
  },
  plugins: [],
}
