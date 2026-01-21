import { createFileRoute } from "@tanstack/react-router";

import { landingClassNames } from "@frontend/ui";

export const Route = createFileRoute("/")({
  component: HomeComponent,
});

function HomeComponent() {
  return (
    <div className={landingClassNames.web.page}>
      <div className={landingClassNames.web.title}>Landing page</div>
    </div>
  );
}
