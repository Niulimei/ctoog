import pkg from './package.json'

export default { ff
  input: 'dist/index.js',
  output: [
    {
      file: pkg['module'],
      format: 'es'
    }
  ]
}
