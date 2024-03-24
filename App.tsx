// import { StatusBar } from "expo-status-bar";
// import { StyleSheet, Text, View } from "react-native";

// export default function App() {
//   return (
//     <View style={styles.container}>
//       <Text>Open up App.tsx to start working on your app!</Text>
//       <StatusBar style="auto" />
//     </View>
//   );
// }

// const styles = StyleSheet.create({
//   container: {
//     flex: 1,
//     backgroundColor: "#fff",
//     alignItems: "center",
//     justifyContent: "center",
//   },
// });

///////////////////////////////////////////////////////
import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import LoginScreen from './screens/LoginScreen';
import ForgotPasswordScreen from './screens/ForgotPasswordScreen';
import SignUpScreen from './screens/SignUpScreen'; 
import CheckEmailScreen from './screens/CheckEmailScreen'; 
import ConfirmationScreen from './screens/ConfirmationScreen'; 
import MainScreen from './screens/MainScreen'; 

export type RootStackParamList = {
  LoginScreen: undefined;
  ForgotPasswordScreen: undefined;
  SignUpScreen: undefined;
  CheckEmailScreen: undefined;
  ConfirmationScreen: undefined;
};

const Stack = createStackNavigator<RootStackParamList>();

const App: React.FC = () => {
  return (
    <NavigationContainer>
      <Stack.Navigator initialRouteName="LoginScreen">
        <Stack.Screen name="LoginScreen" component={LoginScreen} options={{ headerShown: false }} />
        <Stack.Screen name="ForgotPasswordScreen" component={ForgotPasswordScreen} options={{ title: 'Forgot Password' }} />
        <Stack.Screen name="SignUpScreen" component={SignUpScreen} options={{ title: 'Sign Up' }} />
        <Stack.Screen name="CheckEmailScreen" component={CheckEmailScreen} options={{ title: 'Check Email' }} />
        <Stack.Screen name="ConfirmationScreen" component={ConfirmationScreen} options={{ title: 'Confirm' }} />
        <Stack.Screen name="MainScreen" component={MainScreen} options={{ title: 'Main' }} />
        {/* You can add more screens to the navigator as needed */}
      </Stack.Navigator>
    </NavigationContainer>
  );
};

export default App;


// Comments in return function will lead to error: "text should be in..."



//login page 的 LOGIN

//figma setting needs class,major