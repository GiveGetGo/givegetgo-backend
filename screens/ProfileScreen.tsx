import React from 'react';
import { StyleSheet, SafeAreaView, FlatList, View} from 'react-native';
import { Text, Avatar, Divider, Card, Title, Paragraph, Appbar} from 'react-native-paper';
import { Rating } from 'react-native-ratings';
import { useNavigation } from '@react-navigation/native';
import { useSelector } from 'react-redux';
import { useFonts, Montserrat_700Bold_Italic } from '@expo-google-fonts/montserrat';

// Define the types for your navigation stack
type RootStackParamList = {
  SignUpScreen: undefined;
  CheckEmailScreen: undefined;
  LoginScreen: undefined;
};

type UserInfo = {
  name: string;
  email: string;
  bio: string;
  profilePicture: string; // This should be a URI string
  major: string;
  classYear: number;
};

type History = {
  id: string;
  title: string;
  subtitle: string;
  description: string;
  time: string;
  rating?: number; // Optional rating, only for "Match" items
};

const currentUserInfo: UserInfo = { //should be loaded from the database    
  name: 'Gilbert Hsu',
  email: 'xxx@gmail.com',
  bio: 'I am Gilbert the magician.',
  profilePicture: '', 
  major: 'Computer Science',                     
  classYear: 2024,
};

const history_data: History[] = [ //should be loaded from the database //[] because there are more than one data
  {
    id: '1',                             
    title: 'Match',
    subtitle: '$10 Panera Giftcard ',
    description: 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx...',               
    time: '1 hour ago',
    rating: 4,
  },
  {
    id: '2',
    title: 'Request',
    subtitle: '$10 Panera Giftcard ',
    description: 'xxxxxxxxxxxxxxxxx...',
    time: '2 hour ago',
  },
  {
    id: '3',                             
    title: 'Match',
    subtitle: '$10 Panera Giftcard ',
    description: 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx...',               
    time: '1 hour ago',
    rating: 1,
  },
];

// BACKEND DEV SEE HERE: This page could be navigated from either HomeScreen or PostDetailsScreen (both in home stack). Loaded data should differ.

const ProfileScreen: React.FC = () => {   

  const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic });
  
  const use_navigation = useNavigation(); //for Appbar.BackAction

  // get selectedAvatarUri from Redux
  const selectedAvatarUri = useSelector((state: { avatar: { selectedAvatarUri: string } }) => state.avatar.selectedAvatarUri); 
  console.log("selectedAvatarUri from redux:", selectedAvatarUri) 


  // require() could not directly take the url
  interface ImageMap {
      [key: string]: NodeRequire;
  }
  const imageMap: { [key: string]: NodeRequire } = {
      '../assets/avatars/avatar1.png': require('../assets/avatars/avatar1.png'),
      '../assets/avatars/avatar2.png': require('../assets/avatars/avatar2.png'),
      '../assets/avatars/avatar3.png': require('../assets/avatars/avatar3.png'),
      '../assets/avatars/avatar4.png': require('../assets/avatars/avatar4.png'),
      '../assets/avatars/avatar5.png': require('../assets/avatars/avatar5.png'),
      '../assets/avatars/avatar6.png': require('../assets/avatars/avatar6.png'),
      '../assets/avatars/avatar7.png': require('../assets/avatars/avatar7.png'),
      '../assets/avatars/avatar8.png': require('../assets/avatars/avatar8.png'),
      '../assets/avatars/avatar9.png': require('../assets/avatars/avatar9.png'),
  };
  const getProfilePictureSource = (uri: string) => {
      return imageMap[uri] || require(`./profile_icon.jpg`);
  };


  const renderItem = ({ item }: { item: History })  => (
    <Card style={styles.history_card}>
      <Card.Content>
        <View style={styles.titleRow}>
          <Title style={styles.history_title}>{item.title}</Title>
          <Paragraph style={styles.history_time}>{item.time}</Paragraph>
        </View>
        <View style={styles.subtitleRow}>
          <Paragraph style={styles.history_subtitle}>{item.subtitle}</Paragraph>
          {item.title === 'Match' && item.rating !== undefined && (
            <Rating
              type='custom' //so that "ratingBackgroundColor" works
              ratingCount={5}
              imageSize={20}
              readonly
              startingValue={item.rating}
              tintColor="#f6f6f6"
              ratingBackgroundColor='#c8c7c8' 
            />
          )}
        </View>
        <Paragraph style={styles.history_description}>{item.description}</Paragraph>
      </Card.Content>
    </Card>
  );

  return (
    <SafeAreaView  style={styles.container}>
      <View style={styles.headerContainer}>
        <Appbar.BackAction style={styles.backAction} onPress={() => use_navigation.goBack()} />
        <Text style={styles.header}>GiveGetGo</Text>
        <View style={styles.backActionPlaceholder} />
      </View>
      <Avatar.Image source={getProfilePictureSource(selectedAvatarUri)} size={70} style={styles.settings_avatar} /> 
      <Text style={styles.name}>{currentUserInfo.name}</Text>
      <Text style={styles.details}>{`Class of ${currentUserInfo.classYear} • ${currentUserInfo.major}`}</Text>
      <Text style={styles.details_bio}>{currentUserInfo.bio}</Text>
      <Divider style={styles.divider} />
      <FlatList
          data={history_data}
          renderItem={renderItem}
          keyExtractor={(item) => item.id}
          // removeClippedSubviews={true} // If you want to improve performance on large lists, uncomment the line
      />
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginTop: 50,
  },
  headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    justifyContent: 'space-between', // Distributes items evenly
    paddingLeft: 10, // Adds padding to the left of the avatar
    paddingRight: 10, // Adds padding to the right side
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
  settings_avatar: {
    alignSelf: 'center',
    marginTop: 25,
    backgroundColor: '#c7c7c7', // Placeholder color
  },
  name: {
      fontSize: 20,
      fontWeight: 'bold',
      textAlign: 'center',
      marginTop: 8,
      marginBottom: 4,
  },
  details: {
      textAlign: 'center',
      marginBottom: 4,
  },
  details_bio: {
      textAlign: 'center',
      fontStyle: 'italic',
      marginBottom: 6,
  },
  divider: {
      marginVertical: 8,
  },
  history_card: { //page gets longer when there are more contexts
    borderRadius: 15, // Add rounded corners to the card
    marginVertical: 6,
    marginHorizontal: 12,
    elevation: 0, // Adjust for desired shadow depth
    // backgroundColor: '#ffffff', 
    padding: 5, // Add padding inside the card
  },
  titleRow: { //binding title and time together on the same row
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  subtitleRow: { //binding subtitle and rating together on the same row
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 4, 
  },
  history_title: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 0,
  },
  history_subtitle: {
    fontSize: 14, // Adjust the font size as needed
    color: '#000', // Adjust the color as needed
    marginBottom: 0,
  },
  history_time: {
    fontSize: 12,
    color: 'gray', // Adjust the color for the time text
  },
  history_description: {
    fontSize: 14,
    marginBottom: 4, // Adjust spacing to match the image provided
  },
});

export default ProfileScreen;
