const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');

module.exports = (env) => {
  // const env = dotenv.config().parsed;

  // // reduce it to a nice object, the same as before
  // const envKeys = Object.keys(env).reduce((prev, next) => {
  //   prev[`process.env.${next}`] = JSON.stringify(env[next]);
  //   return prev;
  // }, {});
  const devConfig = {
    entry: {
      main: path.resolve(__dirname, './client/source/index.tsx'),
    },
    module: {
      rules: [
        {
          test: /\.(js|jsx)$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader',
            options: {
              presets: ['@babel/preset-react', '@babel/preset-env']
            }
          },
        },
        {
          test: /\.(ts|tsx)$/,
          use: 'ts-loader',
          exclude: /node_modules/,
        },
        {
          test: /\.css$/,
          use: ['style-loader', 'css-loader']
        },
        {
          test: /\.(woff|woff2|ttf|eot|png|jpg|svg|gif)$/i,
          use: ['file-loader']
        }
      ]
    },
    output: {
      filename: '[name].bundle.js',
      path: path.resolve(__dirname, './client/public/build'),
    },
    optimization: {
      splitChunks: {
        chunks: 'all',
      },
    },
    plugins: [
      new HtmlWebpackPlugin({
        title: 'WSO2 API Manager',
        // Load a custom template (lodash by default)
        template: './client/pages/index.html',
        publicPath: '/build',
        templateParameters: { env: env.production ? 'production': 'development'},
      })
    ],
    resolve: {
      alias: {
        assets: path.resolve(__dirname, 'client/source/assets'),
        auth: path.resolve(__dirname, 'client/source/auth'),
        components: path.resolve(__dirname, 'client/source/components'),
        context: path.resolve(__dirname, 'client/source/context'),
        layout: path.resolve(__dirname, 'client/source/layout'),
        'menu-items': path.resolve(__dirname, 'client/source/menu-items'),
        pages: path.resolve(__dirname, 'client/source/pages'),
        routes: path.resolve(__dirname, 'client/source/routes'),
        themes: path.resolve(__dirname, 'client/source/themes'),
        types: path.resolve(__dirname, 'client/source/types'),
        config: path.resolve(__dirname, 'client/source/config.ts'),

        // For the old UIs
        client: path.resolve(__dirname, 'client/src'),
        AppData: path.resolve(__dirname, 'client/source/data/'),
      },
      extensions: ['.tsx', '.ts', '.js', '.jsx'],
    },
    externals: {
      Settings: 'Settings',
      Config: 'OldSettings',
      Themes: 'AppThemes', // Should use long names for preventing global scope JS variable conflicts
      MaterialIcons: 'MaterialIcons',
    },
    devtool: 'source-map',
  }
  if (env.production) {
    return {
      ...devConfig
    }
  } else {
    return {
      mode: 'development',
      ...devConfig
    }
  }
};