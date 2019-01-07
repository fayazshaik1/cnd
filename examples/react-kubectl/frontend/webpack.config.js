const path = require('path');
const webpack = require('webpack');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const pkg = require('./package.json');

const appPath = path.join(__dirname, '/src');
const buildPath = path.resolve(path.join(__dirname, '/build'));

module.exports = {
  context: appPath,
  mode: 'development',
  entry: [
    './index.jsx'
  ],
  output: {
    filename: 'app.[hash].js',
    path: buildPath,
  },
  resolve: {
    extensions: ['.js', '.jsx', '.scss'],
    modules: [
      path.resolve(path.join(__dirname, '/node_modules')),
      path.resolve(appPath)
    ]
  },
  module: {
    rules: [{
      test: /\.jsx?$/,
      exclude: /node_modules/,
      loaders: ['babel-loader'],
    }, {
      test: /\.js?$/,
      exclude: /node_modules/,
      loaders: ['babel-loader'],
    }, {
      test: /\.css$/,
      use: [{
        loader: 'style-loader'
      }, {
        loader: 'css-loader',
        options: {
          includePaths: [appPath]
        }
      }]
    }, 
    {
      test: /\.(scss|sass)$/,
      use: [{
        loader: 'style-loader'
      }, {
        loader: 'css-loader'
      }, {
        loader: 'fast-sass-loader',
        options: {
          includePaths: [appPath]
        }
      }]
    }, { 
      test: /\.(png|jpg|svg)$/, 
      loader: 'url-loader?limit=100000'
    }],
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: './index.html'
    }),
    new CopyWebpackPlugin([
      { from: 'assets/**/*' }
    ]),
    // Enable HMR.
    new webpack.HotModuleReplacementPlugin()
  ],
  devServer: {
    contentBase: buildPath,
    compress: true,
    hot: true,
    inline: true,
    host: '0.0.0.0',
    port: 9000
  }
};