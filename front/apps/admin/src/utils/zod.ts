import type { FormInstance, FormListFieldData } from 'antd'
import type { Rule } from 'antd/es/form'
import type { NamePath } from 'antd/es/form/interface'
import type z from 'zod'
import { Form } from 'antd'
import { useCallback, useMemo } from 'react'

interface UseFormProps<T extends object> {
  form?: FormInstance<T>
  schema?: z.ZodObject
  onSubmit: (data?: T, error?: FormListFieldData) => Promise<boolean> | Promise<void> | boolean | void
  zodValidator?: <T extends Record<string, any>>(props: ZodValidatorProps<T>) => Rule
}

interface UseFormTypes<T> {
  form: FormInstance<T>
  onFinish?: (formData: T) => Promise<boolean | void>
  rules: Rule[]
}

const mapErrorFromZodIssue = (issues: z.core.$ZodIssue[]) =>
  issues.reduce((obj: Record<string, string[]>, issue) => {
    const fieldName = issue.path.join('.')
    if (!obj[fieldName]) {
      obj[fieldName] = [issue.message]
    }
    else {
      obj[fieldName] = [...obj[fieldName], issue.message]
    }
    return obj
  }, {})

interface ZodValidatorProps<T extends Record<string, any>> {
  form: FormInstance<T>
  schema: z.ZodObject
}

export const fieldZodValidator = <T extends object>(props: ZodValidatorProps<T>): Rule => ({
  async validator(rule: any) {
    const values = props.form.getFieldsValue()
    await props.schema.parseAsync(values).catch((e: z.ZodError) => {
      const errorMap = mapErrorFromZodIssue(e.issues)
      const currentFieldError = errorMap[rule.field]
      if (currentFieldError?.length) {
        throw new Error(currentFieldError[0])
      }
    })
  },
})

export const globalZodValidator = <T extends Record<string, any>>(
  props: ZodValidatorProps<T>,
): Rule => ({
  async validator(rule: any) {
    const values = props.form.getFieldsValue()
    const result = await props.schema.safeParseAsync(values)
    const errorMap = result.success ? {} : mapErrorFromZodIssue(result.error.issues)
    const fields = Object.keys(values)
      .filter(key => values[key] !== undefined)
      .map(key => ({
        name: key as NamePath,
        errors: errorMap[key] ?? [],
      }))
    props.form.setFields(fields)
    const currentFieldError = errorMap[rule.field]
    if (currentFieldError?.length) {
      throw new Error(currentFieldError[0])
    }
  },
})

export const defaultZodValidator = globalZodValidator

export function useZodForm<T extends Record<string, any>>({
  form: propsForm,
  schema,
  onSubmit,
  zodValidator = defaultZodValidator,
}: UseFormProps<T>): UseFormTypes<T> {
  const [form] = Form.useForm<T>(propsForm)
  const rules = useMemo(() => (schema ? [zodValidator({ form, schema })] : []), [schema, form, zodValidator])

  const onFinish = useCallback(
    async (formData: T) => {
      if (!schema) {
        return onSubmit(formData)
      }

      try {
        await schema.parseAsync(formData)
        return await onSubmit(formData)
      }
      catch (e: any) {
        const errorMap = mapErrorFromZodIssue(e.issues)
        const _fields = form.getFieldsValue()
        const fields = Object.keys(_fields).map(key => ({
          name: key as NamePath,
          errors: errorMap[key] ?? [],
        }))
        form.setFields(fields)
        return false // 校验失败 → 弹窗不关闭
      }
    },
    [onSubmit, schema, form],
  )

  return {
    form,
    onFinish,
    rules,
  }
}
