import { defineRule } from 'vee-validate';
import { email, required, min, max } from '@vee-validate/rules';

defineRule('email', (value) => {
  if (email(value) && required(value)) {
    return true;
  }

  return 'A valid email address is required';
});

defineRule('password', (value) => {
  if (
    required(value) &&
    min(value, { length: 6 }) &&
    max(value, { length: 30 })
  ) {
    return true;
  }

  return 'Password must be between 6 and 30 characters';
});

defineRule('confirmPassword', (value, [target], ctx) => {
  if (required(value) && value === ctx.form[target]) {
    return true;
  }

  return 'Passwords must match';
});
