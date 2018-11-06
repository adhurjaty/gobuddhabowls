### Installation

This one package will let you use `jquery` and `jquery-ui` *(v.1.12.1, for jQuery1.7+)* in your project. Use `npm install webpack-jquery-ui` instead of installing all dependencies and loaders separately.

The following dependencies will be installed:

| Package name |
|:-------------:|
| jquery |
| jquery-ui |
| css-loader |
| style-loader |
| file-loader |
| webpack |

### Configuration

##### 1. Require the `path` module in your `webpack.config.js` file
```javascript
var path = require('path');

module.exports = {
  //webpack config
}
```

##### 2. Set output folders for images used by jquery-ui
The example of folders arrangement:
```javascript
	output:{
		path: path.join(__dirname,'public/assets/'),
        publicPath:'assets/'
	},
```

##### 3. Set globals
Use the `ProvidePlugin` constructor in the `plugins` object of `webpack.config.js` file to inject `jquery` implicit **globals**:

```javascript
plugins: [
  new webpack.ProvidePlugin({
      $: "jquery",
      jQuery: "jquery",
      "window.jQuery": "jquery'",
      "window.$": "jquery"
  })
]
```

##### 4. Set loaders
Add the loaders to your `webpack.config.js` file to support `jquery-ui`

```javascript
module: {
  rules: [
    {
      test: /\.(jpe?g|png|gif)$/i,
      loader:"file-loader",
      query:{
        name:'[name].[ext]',
        outputPath:'images/'
        //the images will be emmited to public/assets/images/ folder
        //the images will be put in the DOM <style> tag as eg. background: url(assets/images/image.png);
      }
    },
    {
      test: /\.css$/,
      loaders: ["style-loader","css-loader"]
    }
  ]
}
```

### Usage

To load **all** `jquery-ui` **features** with its **basic css theme** use:
```javascript
require('webpack-jquery-ui');
require('webpack-jquery-ui/css');  //ommit, if you don't want to load basic css theme
```

To load **only** `jquery-ui` **interactions** or **widgets** or **effects** [\[list of features\]](http://jqueryui.com/download/) use:
> don't worry about jquery-ui **core files**. All neccessary dependencies are automatically loaded
```javascript
require('webpack-jquery-ui/interactions');
require('webpack-jquery-ui/widgets');
require('webpack-jquery-ui/effects');
```

To load **only particular interactions** [\[list of interactions\]](http://jqueryui.com/download/) use:
> don't worry about jquery-ui **core files**. All neccessary dependencies are automatically loaded with chosen interaction feature
```javascript
require('webpack-jquery-ui/draggable');
require('webpack-jquery-ui/droppable');
require('webpack-jquery-ui/resizable');
require('webpack-jquery-ui/selectable');
require('webpack-jquery-ui/sortable');
```
To load **only particular widgets** [\[list of widgets\]](http://jqueryui.com/download/) use:
> don't worry about jquery-ui **core files**. All neccessary dependencies are automatically loaded with chosen widget
```javascript
require('webpack-jquery-ui/accordion');
require('webpack-jquery-ui/autocomplete');
require('webpack-jquery-ui/button');
require('webpack-jquery-ui/checkboxradio');
require('webpack-jquery-ui/controlgroup');
require('webpack-jquery-ui/datepicker');
require('webpack-jquery-ui/dialog');
require('webpack-jquery-ui/menu');
require('webpack-jquery-ui/progressbar');
require('webpack-jquery-ui/selectmenu');
require('webpack-jquery-ui/slider');
require('webpack-jquery-ui/spinner');
require('webpack-jquery-ui/tabs');
require('webpack-jquery-ui/tooltip');
```
To load **only particular effects** [\[list of effects\]](http://jqueryui.com/download/) use:
> don't worry about jquery-ui **core files**. All neccessary dependencies are automatically loaded with chosen effect
```javascript
require('webpack-jquery-ui/blind-effect');
require('webpack-jquery-ui/bounce-effect');
require('webpack-jquery-ui/clip-effect');
require('webpack-jquery-ui/drop-effect');
require('webpack-jquery-ui/explode-effect');
require('webpack-jquery-ui/fade-effect');
require('webpack-jquery-ui/fold-effect');
require('webpack-jquery-ui/highlight-effect');
require('webpack-jquery-ui/puff-effect');
require('webpack-jquery-ui/pulsate-effect');
require('webpack-jquery-ui/scale-effect');
require('webpack-jquery-ui/shake-effect');
require('webpack-jquery-ui/size-effect');
require('webpack-jquery-ui/slide-effect');
require('webpack-jquery-ui/transfer-effect');
```

You can also set the `entry` property in the `webpack.config.js` file:

```javascript
entry: [
    "webpack-jquery-ui",
    //"webpack-jquery-ui/accordion";
    //"webpack-jquery-ui/widgets";
    //etc.
    "your-entry-point"
]
```

### Links

* [jQuery-UI docs](https://learn.jquery.com/jquery-ui/environments/bower/)
* [jQuery docs](http://jquery.com/download/)

### See also
* [webpack-icons-installer](https://www.npmjs.com/package/webpack-icons-installer)
* [webpack-bootstrap-installer](https://www.npmjs.com/package/webpack-bootstrap-installer)
* [webpack-jquery-ui](https://www.npmjs.com/package/webpack-jquery-ui)
* [webpack-css-loaders](https://www.npmjs.com/package/webpack-css-loaders)
* [webpack-sass-loaders](https://www.npmjs.com/package/webpack-sass-loaders)
* [webpack-babel-installer](https://www.npmjs.com/package/webpack-babel-installer)
* [webpack-karma-jasmine](https://www.npmjs.com/package/webpack-karma-jasmine)