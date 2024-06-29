import { type BaseTextFieldProps, TextField } from "@mui/material";
import { FC, ReactNode, useMemo } from "react";
import { Control, Controller, type UseControllerProps } from "react-hook-form";

interface FormTextFieldProps extends Pick<BaseTextFieldProps, "type" | "required"> {
  name: string;
  control: Control;
  label?: ReactNode;
  rules?: UseControllerProps["rules"];
}
export const FormTextField: FC<FormTextFieldProps> = ({ name, control, label = name, type, required, rules }) => {

  const finalRules = useMemo(() => {
    const _rules = rules || {} as NonNullable<typeof rules>;
    if (required) {
      _rules.required = true;
    }
    return _rules;
  }, [rules, required]);

  return (
    <Controller name={name} control={control}
      rules={finalRules}
      render={({ field: { onChange, value }, fieldState: { error } }) => (
        <TextField
          value={value} onChange={onChange}
          fullWidth variant="outlined"
          label={label} type={type}
          required={required}
          error={!!error} />
      )} />
  );
};