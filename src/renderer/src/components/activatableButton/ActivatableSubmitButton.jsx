import ActivatableButton from './ActivatableButton'

export default function ActivatableSubmitButton(props) {
  return (
    <ActivatableButton
      active={props.active}
      onAction={props.onAction}
      text={'submit'}
      variant={'primary'}
    />
  )
}
