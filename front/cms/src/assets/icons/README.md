# CMS SVG 图标目录

将 SVG 文件放在当前目录或其子目录中，然后直接通过文件名使用全局组件：

```vue
<SvgIcon name="user" size="20" title="用户" />
<SvgIcon name="menu/settings" :size="24" color="#ef4d58" title="设置" />
```

- `name`：相对于本目录的文件路径，可省略 `.svg` 后缀。
- `size`：同时设置宽高，数字按像素处理，默认值为 `1em`。
- `width`、`height`：需要非正方形图标时分别设置。
- `color`：可选；设置后 SVG 作为单色遮罩显示，不设置则保留文件原始颜色。
- `title`：可选；提供无障碍说明。纯装饰图标可以不传。

建议 SVG 文件包含正确的 `viewBox`，并避免在文件中引用外部资源。
