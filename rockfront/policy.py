import reflex as rx

from rockfront.sidebar import sidebar
from rockfront.navbar import navbar


class FormState(rx.State):
    form_data: dict = {}

    def handle_submit(self, form_data: dict):
        """Handle the form submit."""
        self.form_data = form_data



def policy_form():
    return rx.vstack(
        rx.form(
            rx.vstack(
                rx.input(
                    placeholder="First Name",
                    name="first_name",
                ),
                rx.input(
                    placeholder="Last Name",
                    name="last_name",
                ),
                rx.hstack(
                    rx.checkbox("Checked", name="check"),
                    rx.switch("Switched", name="switch"),
                ),
                rx.button("Submit", type="submit"),
            ),
            on_submit=FormState.handle_submit,
            reset_on_submit=True,
        ),
        rx.divider(),
        rx.heading("Results"),
        rx.text(FormState.form_data.to_string()),
    )


def add_policy_page() -> rx.Component:
    return rx.center(
        rx.vstack(
            navbar(),
            rx.hstack(
                sidebar(),
                rx.vstack(
                    policy_form(),
                    align="center",
                )
            )
        )
    )

