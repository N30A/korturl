const copyToClipboard = async (target) => {
  const element = document.querySelector(target);
  if (!element) return;

  const value = element.value;
  if (!value) return;

  await navigator.clipboard.writeText(value);
};

const showBusy = async (target, message) => {
  const element = document.querySelector(target);
  if (!element) return;

  element.setAttribute("aria-busy", "true");

  if (!message) return;

  element.textContent = message;
};
