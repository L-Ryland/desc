import { FormControl, FormControlLabel, type FormGroupProps, FormLabel, Radio, RadioGroup } from "@mui/material";
import type { FC, ReactNode } from "react";
import { Control, Controller } from "react-hook-form";

interface FormRadioGroupProps extends Pick<FormGroupProps, "row"> {
  name: string;
  control: Control;
  label?: ReactNode;
  options?: Array<{ value: unknown, label: string; }>;
  valueAsNumber?: boolean;
  valueAsDate?: boolean;
}
export const FormRadioGroup: FC<FormRadioGroupProps> = ({
  name, control, row, label = name, options = [],
  valueAsDate, valueAsNumber
}) => {
  return (
    <Controller name={name} control={control}
      render={({ field: { value, onChange } }) => (
        <FormControl>
          <FormLabel>{label}</FormLabel>
          <RadioGroup
            name={name} value={value} onChange={e => {
              let result = valueAsNumber ? Number(e.target.value) : valueAsDate ? Date.parse(e.target.value) : e.target.value as any;
              onChange(result);
            }}
            row={row}
          >
            {options.map(({ value, label }, index) => (
              <FormControlLabel key={index} value={value} control={<Radio />} label={label} />
            ))}
          </RadioGroup>
        </FormControl>
      )} />
  );
};
