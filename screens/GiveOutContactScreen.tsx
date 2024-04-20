import React from 'react';
import { View, StyleSheet, SafeAreaView } from 'react-native';
import { Button, Text, Card, Paragraph } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { useFonts, Montserrat_700Bold_Italic } from '@expo-google-fonts/montserrat';

type RootStackParamList = {
  GiveOutContactScreen: undefined;
  NotificationScreen: undefined;
  HomeScreen: undefined;
};

type ScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'GiveOutContactScreen' | 'HomeScreen' 
>;

type NotificationsProps = NativeStackScreenProps<RootStackParamList, 'NotificationScreen'>;

const GiveOutContactScreen: React.FC<NotificationsProps> = ({ navigation }: NotificationsProps) => {

  const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic });

  const use_navigation = useNavigation<ScreenNavigationProp>();

  const GoToHome = () => {
    navigation.navigate('NotificationScreen'); // first let notification stack get back to NotificationScreen
    use_navigation.navigate('HomeScreen'); // then jump to HomeScreen in main stack
  };

  return (
    <SafeAreaView  style={styles.container}> 
      <View style={styles.headerContainer}>
        <View style={styles.backActionPlaceholder} />
        <Text style={styles.header}>GiveGetGo</Text>
        <View style={styles.backActionPlaceholder} />
      </View>
      <Card style={styles.card}>
        <Card.Content>
          <Paragraph style={styles.paragraph}>
            Jimmy Ho’s Contact number is:
          </Paragraph>
          <Paragraph style={styles.paragraph_userinfo}>
            +17650000000
          </Paragraph>
        </Card.Content>
        <Card.Actions style={styles.cardActions}>
          <Button style={styles.button} mode="contained" onPress={GoToHome}>
            Home
          </Button>
        </Card.Actions>
      </Card>
      <Paragraph style={styles.disclaimer}>
        Please be reminded to take note of the following information and maintain the confidentiality of others' privacy.
      </Paragraph>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,                                
    marginTop: 50,
    justifyContent: 'center',
  },
  headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    justifyContent: 'space-between', // Distributes items evenly horizontally
    paddingLeft: 10, 
    paddingRight: 10, 
    position: 'absolute', // So that while setting card to the vertical middle, it still stays at the same place
    top: 0, 
    left: 0,
    right: 2,
    zIndex: 1, // Ensure the headerContainer is above the card
  },
  header: {
    fontSize: 22, // Increase the font size
    fontWeight: '600', // Make the font weight bold
    fontFamily: 'Montserrat_700Bold_Italic',
    textAlign: 'center', // Center the text
    color: '#444444', // Dark gray color
  },
  backActionPlaceholder: {
    width: 48, // This should match the width of the Appbar.BackAction for balance
    height: 52,
  },
  card: { //page gets longer when there are more contexts
    borderRadius: 15, // Add rounded corners to the card
    marginVertical: 6,
    marginHorizontal: 12,
    elevation: 0, // Adjust for desired shadow depth
    // backgroundColor: '#ffffff', 
    padding: 15, // Add padding inside the card
  },
  paragraph: {
    textAlign: 'center',
    fontSize: 16,
    marginBottom: 5,
  },
  paragraph_userinfo: {
    textAlign: 'center',
    fontWeight: 'bold',
    fontStyle: 'italic',
    fontSize: 16,
    marginBottom: 12,
  },
  button: {
    position: 'absolute', 
    left: 110,
    right: 110, //position, left, right together controls the button's length and horizontal location
    alignSelf: 'center', 
  },
  cardActions: {
    justifyContent: 'center', 
    alignItems: 'center',
    padding: 20,
  },
  disclaimer: {
    textAlign: 'center',
    fontSize: 11,
    color: '#888888',
    marginTop: 16,
    paddingHorizontal: 20,
    position: 'absolute', 
    top: 415, 
    left: 0,
    right: 0,
  },
});

export default GiveOutContactScreen;
