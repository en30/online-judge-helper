const CopyWebpackPlugin = require('copy-webpack-plugin');
const path = require('path');
const fs = require('fs');

const entry = fs.readdirSync('src/entry').reduce((a,e) =>
  ({...a, [path.basename(e, '.ts')]:  path.resolve(__dirname, 'src', 'entry', e)}), {})

module.exports = {
  entry,
  module: {
    rules: [
      {
        test: /\.ts$/,
        use: 'ts-loader'
      }
    ]
  },
  resolve: {
    extensions: [ '.ts', '.js' ]
  },
  output: {
    filename: '[name].js',
    path: path.resolve(__dirname, 'dist')
  },
  devtool: 'cheep-module-source-map',
  plugins: [
    new CopyWebpackPlugin([ { from: 'manifest.json' } ])
  ]
};
