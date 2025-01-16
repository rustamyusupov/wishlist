document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('edit');

  form.addEventListener('submit', async event => {
    event.preventDefault();

    const formData = new FormData(form);
    const url = form.getAttribute('action');
    const body = new URLSearchParams([...formData.entries()]).toString();

    const response = await fetch(url, {
      method: event.submitter.id,
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body,
    });

    if (response.ok) {
      window.location.href = '/';
    }
  });
});
