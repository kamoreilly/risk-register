import { Text, View } from "react-native";

import { Container } from "@/components/container";
import { landingClassNames } from "@frontend/ui";

export default function LandingPage() {
  return (
    <Container>
      <View className={landingClassNames.native.page}>
        <Text className={landingClassNames.native.title}>Landing page</Text>
      </View>
    </Container>
  );
}