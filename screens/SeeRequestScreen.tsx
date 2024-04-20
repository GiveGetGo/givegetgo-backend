import React from 'react';
import { View, StyleSheet, TouchableOpacity, SafeAreaView } from 'react-native';
import { Avatar, Button, Text, Card, Title, Paragraph, Appbar } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { useFonts, Montserrat_700Bold_Italic } from '@expo-google-fonts/montserrat';

// Define the types for your navigation stack
type RootStackParamList = {
  SeeRequestScreen: undefined;
  NotificationScreen: undefined;
  NotificationStackProfileScreen: undefined;
  RequestSucceedScreen: undefined;
};

type NotificationsProps = NativeStackScreenProps<RootStackParamList, 'NotificationScreen'>;

const SeeRequestScreen: React.FC<NotificationsProps> = ({ navigation }: NotificationsProps) => {

  const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic });

  const goToProfile = () => {
    navigation.navigate('NotificationStackProfileScreen');
  };

  const takeAction = () => {
    navigation.navigate('RequestSucceedScreen');
  };

  const use_navigation = useNavigation(); //for Appbar.BackAction

  return (
    <SafeAreaView  style={styles.container}> 
      <View style={styles.headerContainer}>
        <Appbar.BackAction style={styles.backAction} onPress={() => use_navigation.goBack()} />
        <Text style={styles.header}>GiveGetGo</Text>
        <View style={styles.backActionPlaceholder} />
        </View>
          <Card style={styles.card}>
            <Card.Content>
              <View style={styles.avatarContainer}>
                <TouchableOpacity onPress={goToProfile}>
                    <Avatar.Image size={70} source={require('./profile_icon.jpg')} />
                </TouchableOpacity>
              </View>
              <Title style={styles.title}>Jimmy Ho</Title>
              <Title style={styles.boldText}>XXXXXX</Title>
              <Paragraph style={styles.paragraph}>XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX</Paragraph>
            </Card.Content>
            <Card.Actions style={styles.cardActions}>
              <Button style={styles.button} mode="contained" onPress={takeAction}>
                Take!
              </Button>
            </Card.Actions>
          </Card>
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
    right: 0,
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
    height: 48,
  },
  backAction: {
    marginLeft: 0 //This means the relative margin, comparing to the container (?)
  },
  card: { //page gets longer when there are more contexts
    borderRadius: 15, // Add rounded corners to the card
    marginVertical: 6,
    marginHorizontal: 12,
    elevation: 0, // Adjust for desired shadow depth
    // backgroundColor: '#ffffff', 
    padding: 15, // Add padding inside the card
    // marginTop: 170,
  },
  avatarContainer: {
    alignItems: 'center',
    justifyContent: 'center',                    
    marginBottom: 0,
    marginTop: -7,
  },
  title: {
    textAlign: 'center',
    marginBottom: -5,
  },
  boldText: {
    textAlign: 'center',
    fontWeight: 'bold',
    fontSize: 16, 
    marginBottom: -2,
  },
  paragraph: {
    textAlign: 'center',
    fontSize: 14,
    marginVertical: 0,
    marginBottom: 12,
  },
  button: {
    // textAlign: 'center',
    // marginBottom: 10,
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
});

export default SeeRequestScreen;
