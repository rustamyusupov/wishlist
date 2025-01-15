document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('edit');

  form.addEventListener('submit', event => {
    event.preventDefault();
    const hiddenInput = form.querySelector('input[name="_method"]');
    hiddenInput.value = event.submitter.id;
    form.submit();
  });
});
