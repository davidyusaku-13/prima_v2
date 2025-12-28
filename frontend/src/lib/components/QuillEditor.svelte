<script>
  import { onMount, onDestroy } from "svelte";
  import Quill from "quill";

  let { value = $bindable(""), onUploadImage } = $props();

  let editorDiv;
  let quill;

  const toolbarOptions = [
    [{ header: [1, 2, 3, false] }],
    ["bold", "italic", "underline", "strike"],
    [{ list: "ordered" }, { list: "bullet" }],
    ["blockquote", "code-block"],
    ["link"],
    ["image"]
  ];

  onMount(() => {
    quill = new Quill(editorDiv, {
      theme: "snow",
      modules: { toolbar: toolbarOptions }
    });

    quill.on("text-change", () => {
      value = quill.root.innerHTML;
    });

    // Custom image handler
    const imageHandler = () => {
      const input = document.createElement("input");
      input.setAttribute("type", "file");
      input.setAttribute("accept", "image/*");
      input.click();
      input.onchange = async () => {
        const file = input.files[0];
        if (file && onUploadImage) {
          const url = await onUploadImage(file);
          if (url) {
            const range = quill.getSelection();
            quill.insertEmbed(range?.index ?? 0, "image", url);
          }
        }
      };
    };

    quill.getModule("toolbar").addHandler("image", imageHandler);
    quill.root.innerHTML = value;
  });

  onDestroy(() => {
    quill = null;
  });

  $effect(() => {
    if (quill && quill.root.innerHTML !== value) {
      quill.root.innerHTML = value;
    }
  });
</script>

<link
  href="https://cdn.quilljs.com/1.3.6/dist/quill.snow.css"
  rel="stylesheet"
/>
<div bind:this={editorDiv} class="quill-editor"></div>

<style>
  :global(.ql-container) {
    font-family: inherit;
    font-size: 1rem;
  }
  :global(.ql-editor) {
    min-height: 300px;
  }
  .quill-editor :global(.ql-toolbar) {
    border-radius: 0.75rem 0.75rem 0 0;
    border-color: #e2e8f0 !important;
  }
  .quill-editor :global(.ql-container) {
    border-radius: 0 0 0.75rem 0.75rem;
    border-color: #e2e8f0 !important;
  }
</style>
